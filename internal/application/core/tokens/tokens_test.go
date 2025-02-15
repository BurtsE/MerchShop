package tokens

import (
	"MerchShop/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateToken(t *testing.T) {
	tests := []struct {
		user  domain.User
		token string
	}{
		{
			user: domain.User{
				ID:       2,
				Username: "username",
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiVXNlcm5hbWUiOiJ1c2VybmFtZSJ9.W1sTtPor_bRtmDzNFqrXemRrMFB2nMx4kUUq0vBW7tU",
		},
		{
			user: domain.User{
				ID:       3,
				Username: "kafka",
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiVXNlcm5hbWUiOiJrYWZrYSJ9.YaW5gmy2QxhLZtEN37H-30oPLGrcBu0DveBhg2iVqr4",
		},
	}
	handler := NewTokenHandler([]byte("secret-key"))
	for _, test := range tests {
		token, err := handler.Create(test.user)
		assert.Nil(t, err)
		assert.Equal(t, test.token, token)
	}
}

func Test_ParseToken(t *testing.T) {
	tests := []struct {
		user  domain.User
		token string
	}{
		{
			user: domain.User{
				ID:       2,
				Username: "username",
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiVXNlcm5hbWUiOiJ1c2VybmFtZSJ9.W1sTtPor_bRtmDzNFqrXemRrMFB2nMx4kUUq0vBW7tU",
		},
		{
			user: domain.User{
				ID:       3,
				Username: "kafka",
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiVXNlcm5hbWUiOiJrYWZrYSJ9.YaW5gmy2QxhLZtEN37H-30oPLGrcBu0DveBhg2iVqr4",
		},
	}
	handler := NewTokenHandler([]byte("secret-key"))
	for _, test := range tests {
		user, err := handler.Parse(test.token)
		assert.Nil(t, err)
		assert.Equal(t, test.user, user)
	}
}

func Test_ParseToken_TokenMalformed(t *testing.T) {
	handler := NewTokenHandler([]byte("secret-key"))
	token := "malformed"
	_, err := handler.Parse(token)
	assert.NotNil(t, err)
}
