package api

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/ports"
)

var _ ports.APIPort = (*Application)(nil)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (a Application) Info(user domain.User) ([]domain.WalletOperation, error) {
	//TODO implement me
	panic("implement me")
}

func (a Application) SendCoin(sender, receiver domain.User, amount int) (domain.WalletOperation, error) {
	//TODO implement me
	panic("implement me")
}

func (a Application) BuyItem(user domain.User, item string) (domain.WalletOperation, error) {
	//TODO implement me
	panic("implement me")
}

func (a Application) Authorize(user domain.User) (string, error) {
	id, err := a.db.CreateUser(user)
	if err != nil {
		return "", err
	}
	user.ID = id
	token, err := a.auth.CreateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a Application) Authenticate(token string) (domain.User, error) {
	tokenUser, err := a.auth.Auth(token)
	if err != nil {
		return domain.User{}, err
	}
	user, err := a.db.GetUser(tokenUser.ID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
