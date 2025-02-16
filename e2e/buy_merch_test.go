package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"io"
	"net/http"
	"testing"
)

func TestBuyItem2(t *testing.T) {
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
	URI := "http://localhost:8080/api/buy/t-shirt"
	// Создаем запрос
	req, err := http.NewRequest(http.MethodGet, URI, nil)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	req.Close = true

	token, err := getToken(t, client, "username")
	t.Log(token, err)
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
	respBody, _ := io.ReadAll(resp.Body)
	t.Log(string(respBody))
}

func getToken(t *testing.T, client *http.Client, username string) (string, error) {
	URI := "http://localhost:8080/api/auth"
	// Создаем запрос
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
    	"username": "%s",
    	"password": "2020"
	}`, username)))
	t.Log(body)
	req, err := http.NewRequest(http.MethodPost, URI, body)
	if err != nil {
		t.Log(body)
		return "", err
	}
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		t.Log(req)
		return "", err
	}
	defer resp.Body.Close()
	t.Log(resp.StatusCode)
	res := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	t.Log(res)
	return res["token"], nil
}
