package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/request"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

func TestWithdrawalRoutes_processWithdrawal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		requestBody        request.Withdrawal
		mockBehavior       func(s *MockWithdrawalUseCase, req request.Withdrawal)
		expectedStatusCode int
	}{
		{
			name: "Success",
			requestBody: request.Withdrawal{
				Order: "12345678903",
				Sum:   100,
			},
			mockBehavior: func(s *MockWithdrawalUseCase, req request.Withdrawal) {
				s.EXPECT().ProcessWithdrawal(gomock.Any(), 1, req.Order, req.Sum).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Order Number",
			requestBody: request.Withdrawal{
				Order: "invalid",
				Sum:   100,
			},
			mockBehavior: func(s *MockWithdrawalUseCase, req request.Withdrawal) {
				s.EXPECT().ProcessWithdrawal(gomock.Any(), 1, req.Order, req.Sum).Return(entity.ErrInvalidOrderNumber)
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Insufficient Funds",
			requestBody: request.Withdrawal{
				Order: "12345678903",
				Sum:   1000,
			},
			mockBehavior: func(s *MockWithdrawalUseCase, req request.Withdrawal) {
				s.EXPECT().ProcessWithdrawal(gomock.Any(), 1, req.Order, req.Sum).Return(entity.ErrInsufficientFunds)
			},
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name: "Internal Server Error",
			requestBody: request.Withdrawal{
				Order: "12345678903",
				Sum:   100,
			},
			mockBehavior: func(s *MockWithdrawalUseCase, req request.Withdrawal) {
				s.EXPECT().ProcessWithdrawal(gomock.Any(), 1, req.Order, req.Sum).Return(errors.New("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWithdrawalUseCase := NewMockWithdrawalUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			mockAuthUseCase := NewMockAuthUseCase(ctrl)

			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			tc.mockBehavior(mockWithdrawalUseCase, tc.requestBody)
			mockAuthUseCase.EXPECT().ValidateToken(gomock.Any()).Return(1, nil)

			router := gin.New()
			protected := router.Group("")
			protected.Use(middleware.AuthMiddleware(mockAuthUseCase))
			newWithdrawalRoutes(protected, mockWithdrawalUseCase, mockLogger)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/balance/withdraw", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer token")
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestWithdrawalRoutes_getUserWithdrawals(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                 string
		mockBehavior         func(s *MockWithdrawalUseCase)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - With Withdrawals",
			mockBehavior: func(s *MockWithdrawalUseCase) {
				tm := time.Now()
				s.EXPECT().GetUserWithdrawals(gomock.Any(), 1).Return([]*entity.Withdrawal{
					{OrderNumber: "123", Sum: 100, ProcessedAt: tm},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Success - No Withdrawals",
			mockBehavior: func(s *MockWithdrawalUseCase) {
				s.EXPECT().GetUserWithdrawals(gomock.Any(), 1).Return([]*entity.Withdrawal{}, nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Internal Server Error",
			mockBehavior: func(s *MockWithdrawalUseCase) {
				s.EXPECT().GetUserWithdrawals(gomock.Any(), 1).Return(nil, errors.New("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWithdrawalUseCase := NewMockWithdrawalUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			mockAuthUseCase := NewMockAuthUseCase(ctrl)

			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			tc.mockBehavior(mockWithdrawalUseCase)
			mockAuthUseCase.EXPECT().ValidateToken(gomock.Any()).Return(1, nil)

			router := gin.New()
			protected := router.Group("")
			protected.Use(middleware.AuthMiddleware(mockAuthUseCase))
			newWithdrawalRoutes(protected, mockWithdrawalUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
			req.Header.Set("Authorization", "Bearer token")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
