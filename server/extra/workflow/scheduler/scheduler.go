package scheduler

import (
	"context"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/types/meta"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"time"
)

// ScheduleResult represents the result of scheduling a stage.
type ScheduleResult struct {
	// UID of the selected worker.
	SuggestedHost string
	// The number of workers the scheduler evaluated the stage against in the filtering
	// phase and beyond.
	EvaluatedWorkers int
	// The number of workers out of the evaluated ones that fit the stage.
	FeasibleWorkers int
}

type Scheduler struct {
	NextStep func() *meta.QueuedStepInfo

	Error func(*meta.QueuedStepInfo, error)

	ScheduleStep func(ctx context.Context, step *meta.Step) (ScheduleResult, error)

	stop chan struct{}

	SchedulingQueue SchedulingQueue

	nextStartWorkerIndex int
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		SchedulingQueue: NewSchedulingQueue(nil),
		stop:            make(chan struct{}),
	}
	s.NextStep = s.nextStep
	return s
}

func (sched *Scheduler) Run(ctx context.Context) {
	sched.SchedulingQueue.Run()

	go parallelizer.JitterUntil(sched.SchedulingOne, time.Second, 0.0, true, sched.stop)

	<-sched.stop
	flog.Info("scheduler stopped")
	sched.SchedulingQueue.Close()
}

func (sched *Scheduler) Shutdown() {
	sched.stop <- struct{}{}
}

func (sched *Scheduler) SchedulingOne() {
	stepInfo := sched.NextStep()

	if stepInfo == nil || stepInfo.Step == nil {
		return
	}

	step := stepInfo.Step
	if sched.skipStepSchedule(step) {
		return
	}

	// todo assume

	// todo bind
}

func (sched *Scheduler) assume() {

}

func (sched *Scheduler) bind() {

}

func (sched *Scheduler) nextStep() *meta.QueuedStepInfo {
	return &meta.QueuedStepInfo{
		StepInfo: &meta.StepInfo{
			Step: &meta.Step{
				Name:              "1",
				UID:               "1",
				WorkerUID:         "",
				ResourceVersion:   "",
				Generation:        0,
				Finalizers:        nil,
				DeletionTimestamp: nil,
				DagUID:            "",
				NodeId:            "",
				Depend:            nil,
				State:             0,
			},
			ParseError: nil,
		},
		Timestamp:               time.Time{},
		Attempts:                0,
		InitialAttemptTimestamp: time.Time{},
		UnschedulablePlugins:    nil,
	}
}

func (sched *Scheduler) skipStepSchedule(step *meta.Step) bool {
	// step is being deleted
	if step.DeletionTimestamp != nil {
		flog.Info("skip step schedule %s", step.UID)
		return true
	}

	return false
}

func (sched *Scheduler) handleSchedulingFailure(ctx context.Context, stepInfo *meta.QueuedStepInfo, err error, reason string, nominatingInfo *meta.NominatingInfo) {
	sched.Error(stepInfo, err)

	//if sched.SchedulingQueue != nil {
	//	sched.SchedulingQueue.AddNominatedStep(stepInfo.StepInfo, nominatingInfo)
	//}

	// todo update store
}
