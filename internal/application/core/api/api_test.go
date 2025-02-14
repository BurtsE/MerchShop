package api

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/ports"
	"context"
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

func TestInfo(t *testing.T) {
	db := new(mockedDb)
	db.On("UserInventory", mock.Anything, mock.Anything).Return(domain.Inventory{}, nil)
	db.On("UserWallet", mock.Anything, mock.Anything).Return([]domain.WalletOperation{}, nil)
	application := NewApplication(db, tokens.NewTokenHandler([]byte("secret_key")))
	_, _, err := application.Info(domain.User{
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
	})
	assert.Nil(t, err)
}
