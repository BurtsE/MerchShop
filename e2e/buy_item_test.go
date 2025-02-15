package e2e

//
//import (
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestBuyItem(t *testing.T) {
//	// Создаем тестовый сервер
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Проверяем метод и путь
//		assert.Equal(t, http.MethodGet, r.Method)
//		assert.Equal(t, "/api/buy/test-item", r.URL.Path)
//
//		// Проверяем заголовок авторизации
//		authHeader := r.Header.Get("Authorization")
//		assert.Equal(t, "Bearer test-token", authHeader)
//
//		// Возвращаем соответствующий статус код
//		w.WriteHeader(http.StatusOK)
//	}))
//	defer ts.Close()
//
//	// Создаем HTTP-клиент
//	client := &http.Client{}
//
//	// Создаем запрос
//	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/buy/test-item", nil)
//	assert.NoError(t, err)
//
//	// Добавляем заголовок авторизации
//	req.Header.Add("Authorization", "Bearer test-token")
//
//	// Выполняем запрос
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	defer resp.Body.Close()
//
//	// Проверяем статус код
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//}
//
//func TestBuyItem_Unauthorized(t *testing.T) {
//	// Создаем тестовый сервер
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Возвращаем статус код 401
//		w.WriteHeader(http.StatusUnauthorized)
//	}))
//	defer ts.Close()
//
//	// Создаем HTTP-клиент
//	client := &http.Client{}
//
//	// Создаем запрос
//	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/buy/test-item", nil)
//	assert.NoError(t, err)
//
//	// Выполняем запрос
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	defer resp.Body.Close()
//
//	// Проверяем статус код
//	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
//}
//
//func TestBuyItem_BadRequest(t *testing.T) {
//	// Создаем тестовый сервер
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Возвращаем статус код 400
//		w.WriteHeader(http.StatusBadRequest)
//	}))
//	defer ts.Close()
//
//	// Создаем HTTP-клиент
//	client := &http.Client{}
//
//	// Создаем запрос
//	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/buy/test-item", nil)
//	assert.NoError(t, err)
//
//	// Выполняем запрос
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	defer resp.Body.Close()
//
//	// Проверяем статус код
//	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
//}
//
//func TestBuyItem_InternalServerError(t *testing.T) {
//	// Создаем тестовый сервер
//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Возвращаем статус код 500
//		w.WriteHeader(http.StatusInternalServerError)
//	}))
//	defer ts.Close()
//
//	// Создаем HTTP-клиент
//	client := &http.Client{}
//
//	// Создаем запрос
//	req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/buy/test-item", nil)
//	assert.NoError(t, err)
//
//	// Выполняем запрос
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	defer resp.Body.Close()
//
//	// Проверяем статус код
//	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
//}
