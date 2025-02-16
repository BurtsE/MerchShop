package e2e

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"net/http"
	"testing"
)

func TestSendCoins(t *testing.T) {
	compose, err := tc.NewDockerCompose("resources/docker-compose.yml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	// Создаем HTTP-клиент
	client := &http.Client{}
	// Создание пользователя для отправки
	receiverName := "receiver"
	_, err = getToken(t, client, receiverName)
	assert.NoError(t, err)
	URI := "http://localhost:8080/api/sendCoin"
	// Создаем запрос
	req, err := http.NewRequest(http.MethodPost, URI, bytes.NewBuffer([]byte(fmt.Sprintf(`{
  		"toUser": "%s",
  		"amount": 30
	}`, receiverName))))
	assert.NoError(t, err)
	if err != nil {
		return
	}
	req.Close = true
	token, err := getToken(t, client, "username")
	assert.NoError(t, err)
	// Добавляем заголовок авторизации
	req.Header.Add("Authorization", "Bearer "+token)

	// Выполняем запрос
	resp, err := client.Do(req)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
