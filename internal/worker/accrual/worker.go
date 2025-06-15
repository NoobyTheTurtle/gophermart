package accrual

import (
	"context"
	"errors"
	"sync"
	"time"

	customErrors "github.com/NoobyTheTurtle/gophermart/internal/entity/errors"
)

type Worker struct {
	id             int
	accrualUseCase AccrualUseCase
	logger         Logger
	maxRetries     int
	orderCh        chan *OrderTask
	stopCh         <-chan struct{}
}

func NewWorker(
	id int,
	accrualUseCase AccrualUseCase,
	logger Logger,
	maxRetries int,
	orderCh chan *OrderTask,
	stopCh <-chan struct{},
) *Worker {
	return &Worker{
		id:             id,
		accrualUseCase: accrualUseCase,
		logger:         logger,
		maxRetries:     maxRetries,
		orderCh:        orderCh,
		stopCh:         stopCh,
	}
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	w.logger.Info("Worker started", "id", w.id)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Worker stopped: context cancelled", "id", w.id)
			return
		case <-w.stopCh:
			w.logger.Info("Worker stopped: stop signal received", "id", w.id)
			return
		case task, ok := <-w.orderCh:
			if !ok {
				w.logger.Info("Worker stopped: order channel closed", "id", w.id)
				return
			}
			w.processOrderTask(ctx, task)
		}
	}
}

func (w *Worker) processOrderTask(ctx context.Context, task *OrderTask) {
	order := task.Order

	w.logger.Info("Worker: processing order", "id", w.id, "order", order.Number, "attempt", task.Attempt)

	err := w.accrualUseCase.ProcessOrder(ctx, order.Number)
	if err != nil {
		var rateLimitErr *customErrors.RateLimitError
		if errors.As(err, &rateLimitErr) {
			w.logger.Info("Worker: rate limit hit for order, retrying after", "id", w.id, "order", order.Number, "retryAfter", rateLimitErr.RetryAfter, "attempt", task.Attempt)
			w.handleRetry(ctx, task, rateLimitErr.RetryAfter)
			return
		}

		w.logger.Error("Worker: failed to process order", "id", w.id, "order", order.Number, "error", err)
		return
	}

	w.logger.Info("Worker: order processed successfully", "id", w.id, "order", order.Number)
}

func (w *Worker) handleRetry(ctx context.Context, task *OrderTask, retryAfter time.Duration) {
	if task.Attempt >= w.maxRetries {
		w.logger.Info("Worker: max retries exceeded for order", "id", w.id, "order", task.Order.Number)
		return
	}

	go func() {
		timer := time.NewTimer(retryAfter)
		defer timer.Stop()

		select {
		case <-timer.C:
			retryTask := &OrderTask{
				Order:   task.Order,
				Attempt: task.Attempt + 1,
			}

			select {
			case w.orderCh <- retryTask:
				w.logger.Info("Worker: order re-queued for retry", "id", w.id, "order", task.Order.Number, "attempt", retryTask.Attempt)
			case <-w.stopCh:
				return
			}
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		}
	}()
}
