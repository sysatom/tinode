package worker

import (
	"context"
	"fmt"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"github.com/tinode/chat/server/extra/utils/queue"
	"github.com/tinode/chat/server/logs"
	"time"
)

type Worker struct {
	Queue *queue.DeltaFIFO

	stop chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		Queue: queue.NewDeltaFIFOWithOptions(queue.DeltaFIFOOptions{
			KeyFunction: nil,
		}),
		stop: make(chan struct{}),
	}
}

func (m *Worker) Run(ctx context.Context) {

	go parallelizer.JitterUntil(m.work, time.Second, 0.0, true, m.stop)

	<-m.stop
	logs.Info.Println("worker stopped")
}

func (m *Worker) Shutdown() {
	m.stop <- struct{}{}
}

func (m *Worker) work() {
	fmt.Println("work", time.Now().UnixMicro())
}
