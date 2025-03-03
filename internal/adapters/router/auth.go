package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

func (r *Router) WithAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(400, gin.H{"errors": "authorization header missing"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		apictx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		user, err := r.app.Authenticate(apictx, token)
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"errors": err.Error()})
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
