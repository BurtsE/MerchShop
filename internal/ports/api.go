package ports

import (
	"MerchShop/internal/application/core/domain"
	"golang.org/x/net/context"
)

type APIPort interface {
	Info(ctx context.Context, user domain.User) (domain.Inventory, []domain.WalletOperation, error)
	SendCoin(ctx context.Context, sender domain.User, receiverName string, amount int) (domain.WalletOperation, error)
	BuyItem(ctx context.Context, user domain.User, item string) error
	Authorize(ctx context.Context, login, password string) (string, error)
	Authenticate(ctx context.Context, token string) (domain.User, error)
}
