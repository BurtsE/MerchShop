package token

import (
	"MerchShop/internal/application/core/domain"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

type TokenHandler struct {
	secretKey []byte
}

func (h TokenHandler) Auth(tokenString string) (domain.User, error) {
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return h.secretKey, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return domain.User{}, nil
}
func (h TokenHandler) CreateToken(user domain.User) (string, error) {
	return "", nil
}
