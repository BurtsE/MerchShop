package api

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/ports"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var _ ports.DBPort = (*mockedDb)(nil)

type mockedDb struct {
	mock.Mock
}

func (d *mockedDb) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	args := d.Called(ctx, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (d *mockedDb) UpdateUser(ctx context.Context, user domain.User) error {
	args := d.Called(ctx, user)
	return args.Error(0)
}

func (d *mockedDb) User(ctx context.Context, userID uint) (domain.User, error) {
	args := d.Called(ctx, userID)
	return args.Get(0).(domain.User), args.Error(1)
}

func (d *mockedDb) UserByName(ctx context.Context, username string) (domain.User, error) {
	args := d.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (d *mockedDb) UserWallet(ctx context.Context, user domain.User) ([]domain.WalletOperation, error) {
	args := d.Called(ctx, user)
	return args.Get(0).([]domain.WalletOperation), args.Error(1)
}

func (d *mockedDb) UserInventory(ctx context.Context, user domain.User) (domain.Inventory, error) {
	args := d.Called(ctx, user)
	return args.Get(0).(domain.Inventory), args.Error(1)
}

func (d *mockedDb) BuyItem(ctx context.Context, user domain.User, item string) (uint, error) {
	args := d.Called(ctx, user, item)
	return args.Get(0).(uint), args.Error(1)
}

func (d *mockedDb) SendCoins(ctx context.Context, from domain.User, to domain.User, amount int) (uint, error) {
	args := d.Called(ctx, from, to, amount)
	return args.Get(0).(uint), args.Error(1)
}

func Test_Info(t *testing.T) {
	user := domain.User{
		ID:           1,
		Username:     "username",
		PasswordHash: "pass_hash",
		Coins:        450,
		Inventory: domain.Inventory{
			{
				Type:     "t-shirt",
				Quantity: 2,
			},
			{
				Type:     "powerbank",
				Quantity: 1,
			},
		},
	}
	inventory := domain.Inventory{{}, {}, {}}
	ops := []domain.WalletOperation{
		{
			ID:     1,
			Sender: user,
			Value:  30,
			Receiver: domain.User{
				ID:           10,
				Username:     "username",
				PasswordHash: "123",
				Coins:        20,
			},
		},
	}
	db := new(mockedDb)
	db.On("UserInventory", mock.Anything, mock.Anything).Return(inventory, nil)
	db.On("UserWallet", mock.Anything, mock.Anything).Return(ops, nil)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	inv, operations, err := application.Info(nil, user)
	assert.Nil(t, err)
	assert.Equal(t, inventory, inv)
	assert.Equal(t, ops, operations)

}

func Test_SendCoin(t *testing.T) {
	db := new(mockedDb)
	sender := domain.User{
		ID:           1,
		Username:     "username",
		PasswordHash: "pass_hash",
		Coins:        450,
		Inventory: domain.Inventory{
			{
				Type:     "t-shirt",
				Quantity: 2,
			},
			{
				Type:     "powerbank",
				Quantity: 1,
			},
		},
	}
	receiver := domain.User{
		ID:           1,
		Username:     "some-user",
		PasswordHash: "pass_hash",
		Coins:        450,
	}
	db.On("UserByName", mock.Anything, mock.Anything).Return(receiver, nil)
	db.On("SendCoins", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint(3), nil)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))

	op, err := application.SendCoin(nil, sender, receiver.Username, 50)
	assert.Nil(t, err)
	assert.Equal(t, uint(3), op.ID)
	assert.Equal(t, sender, op.Sender)
	assert.Equal(t, receiver.Username, op.Receiver.Username)
}

func Test_SendCoin_DBError(t *testing.T) {
	db := new(mockedDb)
	dbErr := fmt.Errorf("database error error")
	db.On("UserByName", mock.Anything, mock.Anything).Return(domain.User{}, dbErr)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, err := application.SendCoin(nil, domain.User{}, "some-user", 50)
	assert.Error(t, err)
}

func Test_SendCoin_DBError2(t *testing.T) {
	db := new(mockedDb)
	dbErr := fmt.Errorf("database error error")
	db.On("UserByName", mock.Anything, mock.Anything).Return(domain.User{}, nil)
	db.On("SendCoins", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint(0), dbErr)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, err := application.SendCoin(nil, domain.User{}, "some-user", 50)
	assert.Error(t, err)
}

func Test_BuyItem(t *testing.T) {
	db := new(mockedDb)
	db.On("BuyItem", mock.Anything, mock.Anything, mock.Anything).Return(uint(3), nil)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	err := application.BuyItem(nil, domain.User{}, "t-shirt")
	assert.Nil(t, err)
}

func Test_BuyItem_DBError(t *testing.T) {
	db := new(mockedDb)
	dbErr := fmt.Errorf("database error")
	db.On("BuyItem", mock.Anything, mock.Anything, mock.Anything).Return(uint(0), dbErr)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	err := application.BuyItem(nil, domain.User{}, "t-shirt")
	assert.Error(t, err)
}

func Test_Authorize(t *testing.T) {
	db := new(mockedDb)
	db.On("CreateUser", mock.Anything, mock.Anything).Return(domain.User{Username: "username"}, nil)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, err := application.Authorize(nil, "username", "password")
	assert.Nil(t, err)
}

func Test_Authorize_DBError(t *testing.T) {
	db := new(mockedDb)
	dbErr := fmt.Errorf("database error")
	db.On("CreateUser", mock.Anything, mock.Anything).Return(domain.User{Username: "username"}, dbErr)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, err := application.Authorize(nil, "username", "password")
	assert.Error(t, err)
}

func Test_Authorize_PasswordError(t *testing.T) {
	db := new(mockedDb)
	longPassword := "password should not contain more than 72 bytes because of hashing algorithms"
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, err := application.Authorize(nil, "username", longPassword)
	assert.Error(t, err)
}

func Test_Authenticate(t *testing.T) {
	tokenHandler := tokens.NewTokenHandler([]byte("secret_key"))
	user := domain.User{
		ID:           3,
		Username:     "username",
		PasswordHash: "password",
		Coins:        430,
	}
	db := new(mockedDb)
	db.On("User", mock.Anything, mock.Anything).Return(user, nil)
	token, err := tokenHandler.Create(user)

	application := NewApplication(db, tokenHandler)
	resultUser, err := application.Authenticate(nil, token)
	assert.Nil(t, err)
	assert.Equal(t, user, resultUser)
}

func Test_Authenticate_MalformedToken(t *testing.T) {
	tokenHandler := tokens.NewTokenHandler([]byte("secret_key"))
	user := domain.User{
		ID:           3,
		Username:     "username",
		PasswordHash: "password",
		Coins:        430,
	}
	db := new(mockedDb)
	db.On("User", mock.Anything, mock.Anything).Return(user, nil)

	application := NewApplication(db, tokenHandler)
	_, err := application.Authenticate(nil, "malformed token")
	assert.Error(t, err)
}

func Test_Authenticate_DBError(t *testing.T) {
	tokenHandler := tokens.NewTokenHandler([]byte("secret_key"))
	user := domain.User{
		ID:           3,
		Username:     "username",
		PasswordHash: "password",
		Coins:        430,
	}
	dbError := fmt.Errorf("database error")
	db := new(mockedDb)
	db.On("User", mock.Anything, mock.Anything).Return(user, dbError)
	token, _ := tokenHandler.Create(user)
	application := NewApplication(db, tokenHandler)
	_, err := application.Authenticate(nil, token)
	assert.Error(t, err)
}
