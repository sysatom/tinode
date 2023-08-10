package scheduler

import (
	"context"
	"github.com/tinode/chat/server/extra/types/meta"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"github.com/tinode/chat/server/logs"
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

func (sched *Scheduler) Run(ctx context.Context) {
	sched.SchedulingQueue.Run()

	go parallelizer.JitterUntilWithContext(ctx, sched.SchedulingOne, 0, 0.0, true)

	<-sched.stop
	logs.Info.Println("scheduler stopped")
	sched.SchedulingQueue.Close()
}

func (sched *Scheduler) Shutdown() {
	sched.stop <- struct{}{}
}

func (sched *Scheduler) SchedulingOne(ctx context.Context) {
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

	logs.Info.Printf("schedule end step %s", step.UID)
}

func (sched *Scheduler) assume() {

}

func (sched *Scheduler) bind() {

}

func (sched *Scheduler) skipStepSchedule(step *meta.Step) bool {
	// step is being deleted
	if step.DeletionTimestamp != nil {
		logs.Info.Printf("skip step schedule %s", step.UID)
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
