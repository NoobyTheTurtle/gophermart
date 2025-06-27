package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/middleware"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

func TestBalanceRoutes_getUserBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                 string
		mockBehavior         func(s *MockBalanceUseCase, l *MockLogger, a *MockAuthUseCase)
		expectedStatusCode   int
		expectedResponseBody string
		authToken            string
	}{
		{
			name: "Success",
			mockBehavior: func(s *MockBalanceUseCase, l *MockLogger, a *MockAuthUseCase) {
				a.EXPECT().ValidateToken("token").Return(1, nil)
				s.EXPECT().GetUserBalance(gomock.Any(), 1).Return(&entity.UserBalance{
					Current:   100.5,
					Withdrawn: 50.25,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"current":100.5,"withdrawn":50.25}`,
			authToken:            "Bearer token",
		},
		{
			name: "Unauthorized",
			mockBehavior: func(s *MockBalanceUseCase, l *MockLogger, a *MockAuthUseCase) {
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"authentication required"}`,
			authToken:            "",
		},
		{
			name: "Internal Server Error",
			mockBehavior: func(s *MockBalanceUseCase, l *MockLogger, a *MockAuthUseCase) {
				a.EXPECT().ValidateToken("token").Return(1, nil)
				s.EXPECT().GetUserBalance(gomock.Any(), 1).Return(nil, errors.New("db error"))
				l.EXPECT().Error("Failed to get user balance", "userID", 1, "error", gomock.Any())
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
			authToken:            "Bearer token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBalanceUseCase := NewMockBalanceUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			mockAuthUseCase := NewMockAuthUseCase(ctrl)

			tc.mockBehavior(mockBalanceUseCase, mockLogger, mockAuthUseCase)

			router := gin.New()
			apiV1 := router.Group("/api/user")
			apiV1.Use(middleware.AuthMiddleware(mockAuthUseCase))
			newBalanceRoutes(apiV1, mockBalanceUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
