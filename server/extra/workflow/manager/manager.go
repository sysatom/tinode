package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/looplab/fsm"
	"github.com/tinode/chat/server/extra/pkg/dag"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"github.com/tinode/chat/server/extra/utils/queue"
	"time"
)

type Manager struct {
	Queue *queue.DeltaFIFO

	stop chan struct{}

	fsm *fsm.FSM
}

func NewManager() *Manager {
	return &Manager{
		Queue: queue.NewDeltaFIFOWithOptions(queue.DeltaFIFOOptions{
			KeyFunction: JobKeyFunc,
		}),
		stop: make(chan struct{}),
		fsm:  NewJobFSM(),
	}
}

func (m *Manager) Run(ctx context.Context) {

	go parallelizer.JitterUntil(m.pushJob, time.Second, 0.0, true, m.stop)

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stop:
			flog.Info("manager stopped")
			return
		default:
			m.popJob()
		}
	}

}

func (m *Manager) Shutdown() {
	m.stop <- struct{}{}
}

func (m *Manager) pushJob() {

	fmt.Println(m.fsm.Current())
	fmt.Println(m.fsm.AvailableTransitions())
	fmt.Println(m.fsm.Event(context.Background(), "run", 1))
	fmt.Println(m.fsm.Current())
	fmt.Println(m.fsm.AvailableTransitions())
	fmt.Println(m.fsm.Event(context.Background(), "success", 1))
	fmt.Println(m.fsm.Current())
	fmt.Println(m.fsm.AvailableTransitions())

	list, err := store.Chatbot.GetJobsByState(model.JobReady)
	if err != nil {
		flog.Error(err)
		return
	}
	for _, job := range list {
		_, exists, err := m.Queue.Get(job)
		if err != nil {
			flog.Error(err)
			continue
		}
		if exists {
			continue
		}

		err = m.Queue.Add(job)
		if err != nil {
			flog.Error(err)
		}
	}
}

func (m *Manager) popJob() {
	_, err := m.Queue.Pop(func(i interface{}) error {
		if d, ok := i.(queue.Deltas); ok {
			for _, delta := range d {
				if delta.Type != queue.Added {
					return nil
				}
				if j, ok := delta.Object.(*model.Job); ok {
					return m.splitDag(j)
				}
			}
		}
		return nil
	})
	if err != nil {
		flog.Error(err)
	}
}

func (m *Manager) splitDag(job *model.Job) error {
	flog.Info("job:%d split dag", job.ID)

	d, err := store.Chatbot.GetDag(int64(job.DagID))
	if err != nil {
		return err
	}
	list, err := dag.TopologySort(d)
	if err != nil {
		return err
	}

	// create steps
	steps := make([]*model.Step, 0, len(list))
	for _, step := range list {
		steps = append(steps, &model.Step{
			UID:    job.UID,
			Topic:  job.Topic,
			JobID:  job.ID,
			Action: model.JSON{"bot": "demo", "action": "start"}, // todo
			Name:   step.Name,
			State:  step.State,
			NodeID: step.NodeId,
			Depend: step.Depend,
		})
	}
	err = store.Chatbot.CreateSteps(steps)
	if err != nil {
		return err
	}

	// update job state
	return store.Chatbot.UpdateJobState(int64(job.ID), model.JobStart)
}

func JobKeyFunc(obj interface{}) (string, error) {
	if j, ok := obj.(*model.Job); ok {
		return fmt.Sprintf("job-%d", j.ID), nil
	}
	return "", errors.New("unknown object")
}

func NewJobFSM() *fsm.FSM {
	f := fsm.NewFSM(
		"ready",
		fsm.Events{
			{Name: "run", Src: []string{"ready"}, Dst: "start"},
			{Name: "success", Src: []string{"start"}, Dst: "finished"},
			{Name: "cancel", Src: []string{"start"}, Dst: "canceled"},
			{Name: "error", Src: []string{"start"}, Dst: "failed"},
		},
		fsm.Callbacks{
			"before_state": func(_ context.Context, e *fsm.Event) {
				fmt.Println("before_state", e)
			},
			"after_state": func(_ context.Context, e *fsm.Event) {
				fmt.Println("after_state", e)
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
