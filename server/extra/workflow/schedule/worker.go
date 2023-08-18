package schedule

import (
	"context"
	"fmt"
	"github.com/tinode/chat/server/extra/pkg/flog"
)

type Worker struct {
	Queue SchedulingQueue

	stop chan struct{}
}

func NewWorker(queue SchedulingQueue) *Worker {
	return &Worker{
		Queue: queue,
		stop:  make(chan struct{}),
	}
}

func (m *Worker) Run(_ context.Context) {

	//go parallelizer.JitterUntil(m.work, time.Second, 0.0, true, m.stop)

	item, err := m.Queue.Pop()
	if err != nil {
		flog.Error(err)
	}
	fmt.Println("worker run", item)
	_ = m.state()

	<-m.stop
	flog.Info("worker stopped")
}

func (m *Worker) Shutdown() {
	m.stop <- struct{}{}
}

func (m *Worker) work() {

}

func (m *Worker) state() error {
	return nil
}
