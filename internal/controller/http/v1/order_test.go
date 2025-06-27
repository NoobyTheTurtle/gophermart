package v1

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

func TestOrderRoutes_uploadOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		requestBody        string
		mockBehavior       func(s *MockOrderUseCase, orderNumber string)
		expectedStatusCode int
	}{
		{
			name:        "Success - Accepted",
			requestBody: "12345678903",
			mockBehavior: func(s *MockOrderUseCase, orderNumber string) {
				s.EXPECT().UploadOrder(gomock.Any(), gomock.Any(), orderNumber).Return(nil)
			},
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:        "Success - OK",
			requestBody: "12345678903",
			mockBehavior: func(s *MockOrderUseCase, orderNumber string) {
				s.EXPECT().UploadOrder(gomock.Any(), gomock.Any(), orderNumber).Return(entity.ErrOrderAlreadyExists)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Unprocessable Entity - Invalid order number",
			requestBody: "123",
			mockBehavior: func(s *MockOrderUseCase, orderNumber string) {
				s.EXPECT().UploadOrder(gomock.Any(), gomock.Any(), orderNumber).Return(entity.ErrInvalidOrderNumber)
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:        "Conflict - Order already exists by another user",
			requestBody: "12345678903",
			mockBehavior: func(s *MockOrderUseCase, orderNumber string) {
				s.EXPECT().UploadOrder(gomock.Any(), gomock.Any(), orderNumber).Return(entity.ErrOrderAlreadyExistsByAnotherUser)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:        "Internal Server Error",
			requestBody: "12345678903",
			mockBehavior: func(s *MockOrderUseCase, orderNumber string) {
				s.EXPECT().UploadOrder(gomock.Any(), gomock.Any(), orderNumber).Return(errors.New("some error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderUseCase := NewMockOrderUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			mockAuthUseCase := NewMockAuthUseCase(ctrl)
			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			tc.mockBehavior(mockOrderUseCase, tc.requestBody)

			router := gin.New()
			protected := router.Group("")
			protected.Use(middleware.AuthMiddleware(mockAuthUseCase))
			newOrderRoutes(protected, mockOrderUseCase, mockLogger)
			mockAuthUseCase.EXPECT().ValidateToken(gomock.Any()).Return(1, nil)

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Authorization", "Bearer token")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}

func TestOrderRoutes_getUserOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                 string
		mockBehavior         func(s *MockOrderUseCase)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success - With orders",
			mockBehavior: func(s *MockOrderUseCase) {
				accrual := 500.0
				s.EXPECT().GetUserOrders(gomock.Any(), 1).Return([]*entity.Order{
					{Number: "123", Status: "PROCESSED", Accrual: accrual, UploadedAt: time.Time{}},
					{Number: "456", Status: "NEW", UploadedAt: time.Time{}},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `[{"number":"123","status":"PROCESSED","accrual":500,"uploaded_at":"0001-01-01T00:00:00Z"},{"number":"456","status":"NEW","uploaded_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "Success - No orders",
			mockBehavior: func(s *MockOrderUseCase) {
				s.EXPECT().GetUserOrders(gomock.Any(), 1).Return([]*entity.Order{}, nil)
			},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Internal Server Error",
			mockBehavior: func(s *MockOrderUseCase) {
				s.EXPECT().GetUserOrders(gomock.Any(), 1).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderUseCase := NewMockOrderUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			mockAuthUseCase := NewMockAuthUseCase(ctrl)
			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			tc.mockBehavior(mockOrderUseCase)
			mockAuthUseCase.EXPECT().ValidateToken(gomock.Any()).Return(1, nil)

			router := gin.New()
			protected := router.Group("")
			protected.Use(middleware.AuthMiddleware(mockAuthUseCase))
			newOrderRoutes(protected, mockOrderUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			req.Header.Set("Authorization", "Bearer token")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
			}
		})
	}
}
