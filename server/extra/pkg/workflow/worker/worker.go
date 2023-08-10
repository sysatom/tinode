package worker

import (
	"context"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"github.com/tinode/chat/server/extra/utils/queue"
	"github.com/tinode/chat/server/logs"
)

type Worker struct {
	Queue queue.DeltaFIFO

	stop chan struct{}
}

func (m *Worker) Run(ctx context.Context) {

	go parallelizer.JitterUntilWithContext(ctx, m.work, 0, 0.0, true)

	<-m.stop
	logs.Info.Println("manager stopped")
}

func (m *Worker) Shutdown() {
	m.stop <- struct{}{}
}

func (m *Worker) work(ctx context.Context) {

}
