package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func main() {
	client := &http.Client{}
	for range 500 {
		sender := RandStringRunes(10)
		receiver := RandStringRunes(10)
		getToken(client, receiver)
		sendCoins(client, sender, receiver)
		time.Sleep(time.Millisecond * 10)
		log.Println("sent coins")
	}

}

func getToken(client *http.Client, username string) (string, error) {
	URI := "http://localhost:8080/api/auth"
	// Создаем запрос
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
    	"username": "%s",
    	"password": "2020"
	}`, username)))
	req, err := http.NewRequest(http.MethodPost, URI, body)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	return res["token"], nil
}

func sendCoins(client *http.Client, senderName, receiverName string) {
	URI := "http://localhost:8080/api/sendCoin"
	// Создаем запрос
	req, err := http.NewRequest(http.MethodPost, URI, bytes.NewBuffer([]byte(fmt.Sprintf(`{
  		"toUser": "%s",
  		"amount": 30
	}`, receiverName))))
	if err != nil {
		log.Println(err)
		return
	}

	token, err := getToken(client, senderName)
	// Добавляем заголовок авторизации
	req.Header.Add("Authorization", "Bearer "+token)

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус код
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(data))
}
