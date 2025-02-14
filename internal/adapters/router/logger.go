package router

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIP := getUserIP(r)
		requestLogger := log.WithFields(log.Fields{"user_ip": userIP})
		next.ServeHTTP(w, r) // Вызов следующего обработчика
		requestLogger.Println(w.Header())
	})
}

func getUserIP(r *http.Request) string {
	// Получаем IP из заголовка X-Forwarded-For (если используется прокси)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0] // Берем первый IP из списка
	}

	// Если заголовок X-Forwarded-For отсутствует, используем RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}
