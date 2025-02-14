package api

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/application/core/tokens"
	"MerchShop/internal/ports"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var _ ports.APIPort = (*Application)(nil)

type Application struct {
	db    ports.DBPort
	token *tokens.TokenHandler
}

func NewApplication(db ports.DBPort, handler *tokens.TokenHandler) *Application {
	app := &Application{db: db, token: handler}
	return app
}

func (a Application) Info(user domain.User) (domain.Inventory, []domain.WalletOperation, error) {
	ctx := context.Background()
	inventory, err := a.db.UserInventory(ctx, user)
	if err != nil {
		return domain.Inventory{}, nil, fmt.Errorf("info getting inventory: %v", err)
	}
	operations, err := a.db.UserWallet(ctx, user)
	if err != nil {
		return domain.Inventory{}, nil, fmt.Errorf("getting operations: %v", err)
	}
	user.Inventory = inventory
	return inventory, operations, nil
}

func (a Application) SendCoin(sender domain.User, receiverName string, amount int) (domain.WalletOperation, error) {
	ctx := context.Background()
	if !sender.Has(amount) {
		return domain.WalletOperation{}, fmt.Errorf("sender does not have enough money")
	}
	receiver, err := a.db.UserByName(ctx, receiverName)
	if err != nil {
		return domain.WalletOperation{}, fmt.Errorf("getting user to send to: %v", err)
	}
	id, err := a.db.SendCoins(context.Background(), sender, receiver, amount)
	if err != nil {
		return domain.WalletOperation{}, fmt.Errorf("sending coins: %v", err)
	}
	return domain.WalletOperation{
		ID:       id,
		Sender:   sender,
		Value:    amount,
		Receiver: receiver,
	}, nil
}

func (a Application) BuyItem(user domain.User, item string) error {
	_, err := a.db.BuyItem(context.Background(), user, item)
	if err != nil {
		return fmt.Errorf("buy item: %v", err)
	}
	return nil
}

func (a Application) Authorize(username, password string) (string, error) {
	if len(password) > 72 {
		return "", fmt.Errorf("password too long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %v", err)
	}
	userCreation := domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	user, err := a.db.CreateUser(context.Background(), userCreation)
	if err != nil {
		return "", err
	}
	token, err := a.token.Create(user)
	if err != nil {
		return "", fmt.Errorf("creating tokens: %v", err)
	}
	return token, nil
}

func (a Application) Authenticate(token string) (domain.User, error) {
	tokenUser, err := a.token.Parse(token)
	if err != nil {
		return domain.User{}, fmt.Errorf("parcing tokens: %v", err)
	}
	user, err := a.db.User(context.Background(), tokenUser.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("getting user: %v", err)
	}
	//err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(tokenUser.PasswordHash))
	//if err != nil {
	//	return domain.User{}, fmt.Errorf("checking password: %v", err)
	//}
	return user, nil
}
