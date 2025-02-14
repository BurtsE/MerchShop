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
		user, err := r.app.Authenticate(token)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		ctx := context.WithValue(req.Context(), "user", user)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	}
}
