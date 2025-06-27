package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/request"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

func TestAuthRoutes_register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                 string
		requestBody          request.Auth
		mockBehavior         func(s *MockAuthUseCase, auth *entity.Auth)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Register(gomock.Any(), gomock.Any()).Return("test-token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"test-token"}`,
		},
		{
			name: "Conflict - User already exists",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", entity.ErrUserAlreadyExists)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"login already taken"}`,
		},
		{
			name: "Internal Server Error",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", errors.New("some internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUseCase := NewMockAuthUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			authData := &entity.Auth{
				Login:    tc.requestBody.Login,
				Password: tc.requestBody.Password,
			}
			tc.mockBehavior(mockAuthUseCase, authData)

			router := gin.New()
			apiV1 := router.Group("/api/user")
			newAuthRoutes(apiV1, mockAuthUseCase, mockLogger)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
		})
	}

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthUseCase := NewMockAuthUseCase(ctrl)
		mockLogger := NewMockLogger(ctrl)

		router := gin.New()
		apiV1 := router.Group("/api/user")
		newAuthRoutes(apiV1, mockAuthUseCase, mockLogger)

		req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer([]byte(`{"login": "test",}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid request format"}`, w.Body.String())
	})
}

func TestAuthRoutes_login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                 string
		requestBody          request.Auth
		mockBehavior         func(s *MockAuthUseCase, auth *entity.Auth)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Login(gomock.Any(), gomock.Any()).Return("test-token", nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"token":"test-token"}`,
		},
		{
			name: "Unauthorized - Invalid Credentials",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Login(gomock.Any(), gomock.Any()).Return("", entity.ErrInvalidCredentials)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid login/password pair"}`,
		},
		{
			name: "Internal Server Error",
			requestBody: request.Auth{
				Login:    "test",
				Password: "password",
			},
			mockBehavior: func(s *MockAuthUseCase, auth *entity.Auth) {
				s.EXPECT().Login(gomock.Any(), gomock.Any()).Return("", errors.New("some internal error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUseCase := NewMockAuthUseCase(ctrl)
			mockLogger := NewMockLogger(ctrl)
			if tc.expectedStatusCode == http.StatusInternalServerError {
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			}

			authData := &entity.Auth{
				Login:    tc.requestBody.Login,
				Password: tc.requestBody.Password,
			}
			tc.mockBehavior(mockAuthUseCase, authData)

			router := gin.New()
			apiV1 := router.Group("/api/user")
			newAuthRoutes(apiV1, mockAuthUseCase, mockLogger)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())

			if tc.expectedStatusCode == http.StatusOK {
				cookie := w.Result().Cookies()[0]
				assert.Equal(t, "token", cookie.Name)
				assert.Equal(t, "test-token", cookie.Value)
			}
		})
	}

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthUseCase := NewMockAuthUseCase(ctrl)
		mockLogger := NewMockLogger(ctrl)

		router := gin.New()
		apiV1 := router.Group("/api/user")
		newAuthRoutes(apiV1, mockAuthUseCase, mockLogger)

		req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer([]byte(`{"login": "test",}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid request format"}`, w.Body.String())
	})
}
