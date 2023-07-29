package main

import (
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	extraStore "github.com/tinode/chat/server/extra/store"
	extraTypes "github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/extra/types/linkit"
	"github.com/tinode/chat/server/logs"
	"github.com/tinode/chat/server/store/types"
	"net/http"
	"sync/atomic"
	"time"
)

// NewExtraSessionStore initializes an extra session store.
func NewExtraSessionStore(lifetime time.Duration) *SessionStore {
	ss := &SessionStore{
		lru:      list.New(),
		lifeTime: lifetime,

		sessCache: make(map[string]*Session),
	}

	//statsRegisterInt("LiveSessions")
	//statsRegisterInt("TotalSessions")

	return ss
}

// queueOut attempts to send a ServerComMessage to a session write loop;
// it fails, if the send buffer is full.
func (s *Session) queueOutExtra(msg *linkit.ServerComMessage) bool {
	if s == nil {
		return true
	}
	if atomic.LoadInt32(&s.terminating) > 0 {
		return true
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logs.Err.Println("s.queueOutExtra: msg marshal failed", s.sid)
		return false
	}

	select {
	case s.send <- data:
	default:
		// Never block here since it may also block the topic's run() goroutine.
		logs.Err.Println("s.queueOutExtra: session's send queue full", s.sid)
		return false
	}
	if s.isMultiplex() {
		s.scheduleClusterWriteLoop()
	}
	return true
}

func (s *Session) readLoopExtra() {
	defer func() {
		s.closeWS()
		s.cleanUp(false)
	}()

	s.ws.SetReadLimit(globals.maxMessageSize)
	_ = s.ws.SetReadDeadline(time.Now().Add(pongWait))
	s.ws.SetPongHandler(func(string) error {
		_ = s.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Read a ClientComMessage
		_, raw, err := s.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				logs.Err.Println("ws: readLoopExtra", s.sid, err)
			}
			return
		}
		statsInc("IncomingMessagesWebsockTotal", 1)
		s.dispatchRawExtra(raw)
	}
}

// Message received, convert bytes to ClientComMessage and dispatch
func (s *Session) dispatchRawExtra(raw []byte) {
	now := types.TimeNow()
	var msg linkit.ClientComMessage

	if atomic.LoadInt32(&s.terminating) > 0 {
		logs.Warn.Println("s.dispatchExtra: message received on a terminating session", s.sid)
		s.queueOut(ErrLocked("", "", now))
		return
	}

	if len(raw) == 1 && raw[0] == 0x31 {
		// 0x31 == '1'. This is a network probe message. Respond with a '0':
		s.queueOutBytes([]byte{0x30})
		return
	}

	toLog := raw
	truncated := ""
	if len(raw) > 512 {
		toLog = raw[:512]
		truncated = "<...>"
	}
	logs.Info.Printf("in: '%s%s' sid='%s' uid='%s'", toLog, truncated, s.sid, s.uid)

	if err := json.Unmarshal(raw, &msg); err != nil {
		// Malformed message
		logs.Warn.Println("s.dispatchExtra", err, s.sid)
		s.queueOut(ErrMalformed("", "", now))
		return
	}

	s.dispatchExtra(&msg)
}

func (s *Session) dispatchExtra(_ *linkit.ClientComMessage) {
	// pull
	instructs, err := extraStore.Chatbot.ListInstruct(s.uid, false)
	if err != nil {
		s.queueOutExtra(ErrMessage(400, err.Error()))
		return
	}
	var instruct []map[string]interface{}
	instruct = []map[string]interface{}{}
	for _, item := range instructs {
		instruct = append(instruct, map[string]interface{}{
			"no":        item.No,
			"bot":       item.Bot,
			"flag":      item.Flag,
			"content":   item.Content,
			"expire_at": item.ExpireAt,
		})
	}

	s.queueOutExtra(&linkit.ServerComMessage{
		Code: http.StatusOK,
		Data: instruct,
	})
}

// ErrMessage error message with code.
func ErrMessage(code int, message string) *linkit.ServerComMessage {
	return &linkit.ServerComMessage{
		Code:    code,
		Message: message,
	}
}

// Get API key from an HTTP request.
func getAccessToken(req *http.Request) string {
	// Check header.
	apikey := req.Header.Get("X-AccessToken")
	if apikey != "" {
		return apikey
	}

	// Check URL query parameters.
	apikey = req.URL.Query().Get("accessToken")
	if apikey != "" {
		return apikey
	}

	// Check form values.
	apikey = req.FormValue("accessToken")
	if apikey != "" {
		return apikey
	}

	// Check cookies.
	if c, err := req.Cookie("accessToken"); err == nil {
		apikey = c.Value
	}

	return apikey
}

func checkAccessToken(accessToken string) (uid types.Uid, isValid bool) {
	p, err := extraStore.Chatbot.ParameterGet(accessToken)
	if err != nil {
		return
	}
	if p.ID <= 0 || p.IsExpired() {
		return
	}

	u, _ := extraTypes.KV(p.Params).String("uid")
	uid = types.ParseUserId(u)
	if uid.IsZero() {
		return
	}
	isValid = true
	return
}
