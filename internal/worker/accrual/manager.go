package accrual

import (
	"context"
	"sync"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type OrderTask struct {
	Order   *entity.Order
	Attempt int
}

type AccrualManager struct {
	accrualUseCase  AccrualUseCase
	logger          Logger
	workerCount     int
	processInterval time.Duration
	maxRetries      int

	orderCh chan *OrderTask
	stopCh  chan struct{}
	wg      sync.WaitGroup

	coordinator *Coordinator
	workers     []*Worker
}

func NewAccrualManager(accrualUseCase AccrualUseCase, logger Logger, workerCount int, processInterval int) *AccrualManager {
	orderCh := make(chan *OrderTask, workerCount*2)
	stopCh := make(chan struct{})

	am := &AccrualManager{
		accrualUseCase:  accrualUseCase,
		logger:          logger,
		workerCount:     workerCount,
		processInterval: time.Duration(processInterval) * time.Second,
		maxRetries:      3,
		orderCh:         orderCh,
		stopCh:          stopCh,
	}

	am.coordinator = NewCoordinator(
		accrualUseCase,
		logger,
		am.processInterval,
		orderCh,
		stopCh,
	)

	am.workers = make([]*Worker, workerCount)
	for i := range workerCount {
		am.workers[i] = NewWorker(
			i,
			accrualUseCase,
			logger,
			am.maxRetries,
			orderCh,
			stopCh,
		)
	}

	return am
}

func (m *AccrualManager) Start(ctx context.Context) {
	m.logger.Info("Starting accrual manager with coordinator and workers", "count", m.workerCount)

	m.wg.Add(1)
	go m.coordinator.Start(ctx, &m.wg)

	for _, worker := range m.workers {
		m.wg.Add(1)
		go worker.Start(ctx, &m.wg)
	}
}

func (m *AccrualManager) Stop() {
	m.logger.Info("Stopping accrual manager...")
	close(m.stopCh)
	m.wg.Wait()
	close(m.orderCh)
	m.logger.Info("Accrual manager stopped")
}
