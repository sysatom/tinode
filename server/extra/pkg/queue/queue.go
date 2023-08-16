package queue

import (
	"github.com/adjust/rmq/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/types"
	"github.com/tinode/chat/server/logs"
	serverTypes "github.com/tinode/chat/server/store/types"
	"os"
	"time"
)

const (
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
)

var connection rmq.Connection

func init() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	if addr == "" || password == "" {
		panic("redis config error")
	}

	errChan := make(chan error, 10)
	go logErrors(errChan)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	var err error
	connection, err = rmq.OpenConnectionWithRedisClient("consumer", client, errChan)
	if err != nil {
		panic(err)
	}
}

var messageQueue rmq.Queue

func InitMessageQueue(consumer rmq.Consumer) {
	var err error
	messageQueue, err = connection.OpenQueue("messages")
	if err != nil {
		panic(err)
	}

	if err = messageQueue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		panic(err)
	}

	if _, err = messageQueue.AddConsumer("message", consumer); err != nil {
		panic(err)
	}
}

func Shutdown() {
	<-messageQueue.StopConsuming()
	flog.Info("message queue stopped")
}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				logs.Err.Println("heartbeat error (limit): ", err)
			} else {
				logs.Err.Println("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			logs.Err.Println("consume error: ", err)
		case *rmq.DeliveryError:
			logs.Err.Println("delivery error: ", err.Delivery, err)
		default:
			logs.Err.Println("other error: ", err)
		}
		flog.Error(err)
	}
}

func AsyncMessage(rcptTo, original string, msg types.MsgPayload) error {
	botUid := serverTypes.ParseUserId(original)
	qp, err := types.ConvertQueuePayload(rcptTo, botUid.UserId(), msg)
	if err != nil {
		return nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	payload, err := json.Marshal(qp)
	if err != nil {
		return nil
	}
	return messageQueue.PublishBytes(payload)
}

func Stats() (string, error) {
	queues, err := connection.GetOpenQueues()
	if err != nil {
		return "", err
	}

	stats, err := connection.CollectStats(queues)
	if err != nil {
		return "", err
	}

	return stats.GetHtml("", ""), nil
}
