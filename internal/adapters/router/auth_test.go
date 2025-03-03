package router

import (
	"MerchShop/internal/application/core/domain"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockApp представляет собой mock-объект для тестирования.
type MockApp struct {
	authenticateFunc func(ctx context.Context, token string) (domain.User, error)
}

func (m MockApp) Info(ctx context.Context, user domain.User) (domain.Inventory, []domain.WalletOperation, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockApp) SendCoin(ctx context.Context, sender domain.User, receiverName string, amount int) (domain.WalletOperation, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockApp) BuyItem(ctx context.Context, user domain.User, item string) error {
	//TODO implement me
	panic("implement me")
}

func (m MockApp) Authorize(ctx context.Context, login, password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockApp) Authenticate(ctx context.Context, token string) (domain.User, error) {
	return m.authenticateFunc(ctx, token)
}

// TestWithAuth проверяет различные сценарии работы middleware WithAuth.
func TestWithAuth(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		authenticateFn func(ctx context.Context, token string) (domain.User, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			authenticateFn: nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"errors\":\"authorization header missing\"}",
		},
		{
			name:       "Invalid Token",
			authHeader: "Bearer invalid-token",
			authenticateFn: func(ctx context.Context, token string) (domain.User, error) {
				return domain.User{}, errors.New("invalid token")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "{\"errors\":\"invalid token\"}",
		},
		{
			name:           "Valid Token",
			authHeader:     "Bearer valid-token",
			authenticateFn: func(ctx context.Context, token string) (domain.User, error) { return domain.User{}, nil },
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
	}
	rtr := gin.Default()
	router := &Router{
		srv: http.Server{
			Handler: rtr,
		},
	}
	rtr.Use(router.WithAuth())
	// Создаем тестовый обработчик, который будет вызван после WithAuth
	nextHandler := func(ctx *gin.Context) {
		// Проверяем, что пользователь добавлен в контекст
		user, _ := ctx.Get("user")
		if user == nil {
			t.Error("User not found in context")
		}
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write([]byte("OK"))
	}
	rtr.GET("/", nextHandler)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем mock-объект для приложения
			mockApp := &MockApp{
				authenticateFunc: tt.authenticateFn,
			}
			router.app = mockApp

			// Создаем тестовый HTTP-запрос
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Создаем ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()
			// Создаем тестовый обработчик, который будет вызван после WithAuth
			rtr.ServeHTTP(rr, req)

			// Проверяем статус код
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Проверяем тело ответа
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
