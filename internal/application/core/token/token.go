package token

import (
	"MerchShop/internal/application/core/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = fmt.Errorf("token is invalid")

type TokenHandler struct {
	secretKey []byte
}

func NewTokenHandler(secretKey []byte) *TokenHandler {
	return &TokenHandler{secretKey: secretKey}
}

func (h TokenHandler) Auth(tokenString string) (domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return h.secretKey, nil
	})
	if err != nil {
		return domain.User{}, err
	}

	var (
		claims jwt.MapClaims
		ok     bool
		user   domain.User
	)
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if user.ID, ok = claims["ID"].(uint); !ok {
		return domain.User{}, ErrInvalidToken
	}
	if user.ID, ok = claims["Name"].(uint); !ok {
		return domain.User{}, ErrInvalidToken
	}
	return user, nil
}
func (h TokenHandler) CreateToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   user.ID,
		"Name": user.Name,
	})
	return token.SignedString(h.secretKey)
}
