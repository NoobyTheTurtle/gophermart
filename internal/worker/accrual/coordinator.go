package accrual

import (
	"context"
	"sync"
	"time"
)

type Coordinator struct {
	accrualUseCase  AccrualUseCase
	logger          Logger
	processInterval time.Duration
	orderCh         chan<- *OrderTask
	stopCh          <-chan struct{}
}

func NewCoordinator(
	accrualUseCase AccrualUseCase,
	logger Logger,
	processInterval time.Duration,
	orderCh chan<- *OrderTask,
	stopCh <-chan struct{},
) *Coordinator {
	return &Coordinator{
		accrualUseCase:  accrualUseCase,
		logger:          logger,
		processInterval: processInterval,
		orderCh:         orderCh,
		stopCh:          stopCh,
	}
}

func (c *Coordinator) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	c.logger.Info("Coordinator started")
	ticker := time.NewTicker(c.processInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Coordinator stopped: context cancelled")
			return
		case <-c.stopCh:
			c.logger.Info("Coordinator stopped: stop signal received")
			return
		case <-ticker.C:
			c.fetchAndDistributeOrders(ctx)
		}
	}
}

func (c *Coordinator) fetchAndDistributeOrders(ctx context.Context) {
	orders, err := c.accrualUseCase.GetOrdersForProcessing(ctx)
	if err != nil {
		c.logger.Error("Coordinator: failed to get orders for processing", "error", err)
		return
	}

	if len(orders) == 0 {
		return
	}

	c.logger.Info("Coordinator: distributing orders to workers", "count", len(orders))

	for _, order := range orders {
		task := &OrderTask{
			Order:   order,
			Attempt: 1,
		}

		select {
		case c.orderCh <- task:
		case <-c.stopCh:
			c.logger.Info("Coordinator: stopping during order distribution")
			return
		default:
			c.logger.Info("Coordinator: worker queue is full, skipping order, will retry on next iteration", "order", order.Number)
		}
	}
}
