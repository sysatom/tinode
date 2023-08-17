package scheduler

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
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

func (m *Worker) Run(ctx context.Context) {

	//go parallelizer.JitterUntil(m.work, time.Second, 0.0, true, m.stop)

	item, err := m.Queue.Pop()
	if err != nil {
		flog.Error(err)
	}
	fmt.Println(item)
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

func NewStepFSM() *fsm.FSM {
	f := fsm.NewFSM(
		"created",
		fsm.Events{
			{Name: "bind", Src: []string{"created"}, Dst: "ready"},
			{Name: "run", Src: []string{"ready"}, Dst: "running"},
			{Name: "success", Src: []string{"running"}, Dst: "finished"},
			{Name: "error", Src: []string{"running"}, Dst: "failed"},
			{Name: "cancel", Src: []string{"running"}, Dst: "canceled"},
			{Name: "skip", Src: []string{"running"}, Dst: "skipped"},
		},
		fsm.Callbacks{
			"before_state": func(_ context.Context, e *fsm.Event) {
				fmt.Println("before_state")
			},
			"after_state": func(_ context.Context, e *fsm.Event) {
				fmt.Println("after_state")
			},
		},
	)

	s, err := fsm.VisualizeWithType(f, fsm.MERMAID)
	if err != nil {
		flog.Error(err)
	}
	fmt.Println(s)

	return f
}
