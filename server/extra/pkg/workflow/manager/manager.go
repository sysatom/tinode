package manager

import (
	"context"
	"github.com/tinode/chat/server/extra/utils/parallelizer"
	"github.com/tinode/chat/server/extra/utils/queue"
	"github.com/tinode/chat/server/logs"
)

type Manager struct {
	Queue queue.DeltaFIFO

	stop chan struct{}
}

func (m *Manager) Run(ctx context.Context) {

	go parallelizer.JitterUntilWithContext(ctx, m.manager, 0, 0.0, true)

	<-m.stop
	logs.Info.Println("manager stopped")
}

func (m *Manager) Shutdown() {
	m.stop <- struct{}{}
}

func (m *Manager) manager(ctx context.Context) {

}
