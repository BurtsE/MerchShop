package ports

import "MerchShop/internal/application/core/domain"

type APIPort interface {
	Info(user domain.User) (domain.Inventory, []domain.WalletOperation, error)
	SendCoin(sender domain.User, receiverName string, amount int) (domain.WalletOperation, error)
	BuyItem(user domain.User, item string) error
	Authorize(login, password string) (string, error)
	Authenticate(token string) (domain.User, error)
}
