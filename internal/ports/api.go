package ports

import "MerchShop/internal/application/core/domain"

type APIPort interface {
	Info(user domain.User) ([]domain.WalletOperation, error)
	SendCoin(sender, receiver domain.User, amount int) (domain.WalletOperation, error)
	BuyItem(user domain.User, item string) (domain.WalletOperation, error)
	Authorize(user domain.User) (string, error)
	Authenticate(token string) (domain.User, error)
}
