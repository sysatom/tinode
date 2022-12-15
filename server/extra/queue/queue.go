package queue

import (
	"encoding/json"
	"github.com/adjust/rmq/v5"
	"github.com/go-redis/redis/v8"
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

var MessageQueue rmq.Queue

func InitMessageQueue(consumer rmq.Consumer) {
	var err error
	MessageQueue, err = connection.OpenQueue("messages")
	if err != nil {
		panic(err)
	}

	if err = MessageQueue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		panic(err)
	}

	if _, err = MessageQueue.AddConsumer("message", consumer); err != nil {
		panic(err)
	}
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
	}
}

func AsyncMessage(rcptTo, original string, msg types.MsgPayload) error {
	botUid := serverTypes.ParseUserId(original)
	qp, err := types.ConvertQueuePayload(rcptTo, botUid.UserId(), msg)
	if err != nil {
		return nil
	}
	payload, err := json.Marshal(qp)
	if err != nil {
		return nil
	}
	return MessageQueue.PublishBytes(payload)
}
