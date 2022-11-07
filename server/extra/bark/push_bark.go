package bark

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tinode/chat/server/drafty"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/push"
	serverStore "github.com/tinode/chat/server/store"
	"github.com/tinode/chat/server/store/types"
	"io"
	"net/http"
	"os"
	"time"
)

var handler barkPush

// How much to buffer the input channel.
const defaultBuffer = 32

// BarkDeviceKey store device key
const BarkDeviceKey = "bark_device_key"

type barkPush struct {
	initialized bool
	input       chan *push.Receipt
	channel     chan *push.ChannelReq
	stop        chan bool
}

type configType struct {
	Enabled   bool   `json:"enabled"`
	Buffer    int    `json:"buffer"`
	Sound     string `json:"sound"`
	Icon      string `json:"icon"`
	ApiUrl    string `json:"api_url"`
	DeviceKey string `json:"device_key"`
}

// Init initializes the handler
func (barkPush) Init(jsonconf string) error {
	// Check if the handler is already initialized
	if handler.initialized {
		return errors.New("already initialized")
	}

	var config configType
	if err := json.Unmarshal([]byte(jsonconf), &config); err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	handler.initialized = true

	if !config.Enabled {
		return nil
	}

	if config.Buffer <= 0 {
		config.Buffer = defaultBuffer
	}

	handler.input = make(chan *push.Receipt, config.Buffer)
	handler.channel = make(chan *push.ChannelReq, config.Buffer)
	handler.stop = make(chan bool, 1)

	go func() {
		for {
			select {
			case msg := <-handler.input:
				sendPushes(&config, msg)
			case msg := <-handler.channel:
				fmt.Fprintln(os.Stdout, msg)
			case <-handler.stop:
				return
			}
		}
	}()

	return nil
}

// IsReady checks if the handler is initialized.
func (barkPush) IsReady() bool {
	return handler.input != nil
}

// Push returns a channel that the server will use to send messages to.
// If the adapter blocks, the message will be dropped.
func (barkPush) Push() chan<- *push.Receipt {
	return handler.input
}

// Channel returns a channel that caller can use to subscribe/unsubscribe devices to channels (FCM topics).
// If the adapter blocks, the message will be dropped.
func (barkPush) Channel() chan<- *push.ChannelReq {
	return handler.channel
}

// Stop terminates the handler's worker and stops sending pushes.
func (barkPush) Stop() {
	handler.stop <- true
}

func sendPushes(config *configType, rcpt *push.Receipt) {
	body, err := drafty.PlainText(rcpt.Payload.Content)
	if err != nil {
		logs.Err.Println(err)
		return
	}

	for uid := range rcpt.To {
		if uid.UserId() == rcpt.Payload.From {
			continue
		}

		v, err := store.Chatbot.ConfigGet(uid, "", BarkDeviceKey)
		if err != nil {
			logs.Err.Println(err)
			continue
		}
		config.DeviceKey, _ = v.String("value")

		from := types.ParseUserId(rcpt.Payload.From)
		if from.IsZero() {
			from = types.ParseUid(rcpt.Payload.From)
		}
		fromUser, err := serverStore.Users.Get(from)
		if err != nil {
			logs.Err.Println(err)
			continue
		}
		if fromUser != nil && fromUser.Public != nil {
			if public, ok := fromUser.Public.(map[string]interface{}); ok {
				name := public["fn"].(string)
				body = fmt.Sprintf("[%s] %s", name, body)
			}
		}

		err = postMessage(config, "", body, rcpt.Payload.Topic)
		if err != nil {
			logs.Err.Println(err)
			return
		}
	}
}

func postMessage(config *configType, title, body, group string) error {
	if config.DeviceKey == "" {
		return errors.New("device key empty")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	m := map[string]interface{}{
		"title":      title,
		"body":       body,
		"device_key": config.DeviceKey,
		"sound":      config.Sound,
		"icon":       config.Icon,
		"group":      group,
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	data := &bytes.Buffer{}
	_, err = data.Write(j)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.ApiUrl, data)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logs.Info.Println(string(b))

	return nil
}

func init() {
	push.Register("bark", &handler)
}
