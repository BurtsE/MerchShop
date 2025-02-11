package router

import "MerchShop/internal/application/core/domain"

type UserData struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoin struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoin struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedCoin `json:"received"`
	Sent     []SentCoin     `json:"sent"`
}

func ConvertToUserData(user domain.User, operations []domain.WalletOperation) UserData {
	inventory := make([]InventoryItem, len(user.Inventory))
	for i, item := range user.Inventory {
		inventory[i] = InventoryItem{
			Type:     item.Type,
			Quantity: item.Quantity,
		}
	}

	var received []ReceivedCoin
	var sent []SentCoin

	for _, op := range operations {
		switch user.ID {
		case op.Receiver.ID:
			received = append(received, ReceivedCoin{
				FromUser: op.Sender.Name,
				Amount:   op.Value,
			})
		case op.Sender.ID:
			sent = append(sent, SentCoin{
				ToUser: op.Receiver.Name,
				Amount: op.Value,
			})
		}
	}

	// Создаем и возвращаем UserData
	return UserData{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}
}
