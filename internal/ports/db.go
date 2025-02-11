package ports

import "MerchShop/internal/application/core/domain"

type DBPort interface {
	CreateUser(user domain.User) (uint, error)
	GetUser(userId uint) (domain.User, error)
	UpdateUser(user domain.User) error
	BuyItem(user domain.User, item string) (uint, error)
	SendCoins(from domain.User, to domain.User, amount int) error
}
