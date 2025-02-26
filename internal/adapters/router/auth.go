package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (r *Router) WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("authorization header missing"))
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		user, err := r.app.Authenticate(ctx, token)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		contextWithUserValue := context.WithValue(req.Context(), "user", user)
		req = req.WithContext(contextWithUserValue)
		next.ServeHTTP(w, req)
	}
}
