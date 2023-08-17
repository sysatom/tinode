package scheduler

import (
	"context"
	"github.com/tinode/chat/server/extra/pkg/flog"
	"github.com/tinode/chat/server/extra/store"
	"github.com/tinode/chat/server/extra/store/model"
	"github.com/tinode/chat/server/extra/types/meta"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"time"
)

type Scheduler struct {
	NextStep func() *meta.QueuedStepInfo

	stop chan struct{}

	SchedulingQueue SchedulingQueue

	nextStartWorkerIndex int
}

func NewScheduler(queue SchedulingQueue) *Scheduler {
	s := &Scheduler{
		SchedulingQueue: queue,
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

	err := sched.SchedulingQueue.Add(step)
	if err != nil {
		flog.Error(err)
	}
}

func (sched *Scheduler) nextStep() *meta.QueuedStepInfo {
	readyStep, err := store.Chatbot.GetStepByState(model.StepReady)
	if err != nil {
		flog.Error(err)
		return nil
	}

	return &meta.QueuedStepInfo{
		StepInfo: &meta.StepInfo{
			Step: &meta.Step{
				Name:            readyStep.Name,
				UID:             "",
				WorkerUID:       "",
				ResourceVersion: "",
				Generation:      0,
				Finalizers:      nil,
				JobId:           readyStep.JobID,
				NodeId:          readyStep.NodeID,
				Depend:          readyStep.Depend,
				State:           readyStep.State,
			},
			ParseError: nil,
		},
		Timestamp:               readyStep.CreatedAt,
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

	//if sched.SchedulingQueue != nil {
	//	sched.SchedulingQueue.AddNominatedStep(stepInfo.StepInfo, nominatingInfo)
	//}

	// todo update store
}
