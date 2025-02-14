package ports

import (
	"MerchShop/internal/application/core/domain"
	"context"
)

type DBPort interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	User(ctx context.Context, userID uint) (domain.User, error)
	UserByName(ctx context.Context, username string) (domain.User, error)
	UserWallet(ctx context.Context, user domain.User) ([]domain.WalletOperation, error)
	UserInventory(ctx context.Context, user domain.User) (domain.Inventory, error)
	BuyItem(ctx context.Context, user domain.User, item string) (uint, error)
	SendCoins(ctx context.Context, from domain.User, to domain.User, amount int) (uint, error)
}
