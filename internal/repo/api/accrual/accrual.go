package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/entity/errors"
)

type AccrualAPI struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *AccrualAPI {
	return &AccrualAPI{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AccrualAPI) GetOrderAccrual(ctx context.Context, orderNumber string) (*entity.Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("repo - api - accrual - GetOrderAccrual: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo - api - accrual - GetOrderAccrual: failed to make request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var accrual entity.Accrual
		if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
			return nil, fmt.Errorf("repo - api - accrual - GetOrderAccrual: failed to decode response: %w", err)
		}
		return &accrual, nil
	case http.StatusNoContent:
		return nil, entity.ErrOrderNotFound
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				return nil, &errors.RateLimitError{RetryAfter: time.Duration(seconds) * time.Second}
			}
		}
		return nil, &errors.RateLimitError{RetryAfter: 60 * time.Second}
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("repo - api - accrual - GetOrderAccrual: accrual system internal error")
	default:
		return nil, fmt.Errorf("repo - api - accrual - GetOrderAccrual: unexpected status code: %d", resp.StatusCode)
	}
}
