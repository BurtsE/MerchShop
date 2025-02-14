package domain

type User struct {
	ID           uint
	Username     string
	PasswordHash string
	Coins        int
	Inventory    Inventory
}
type Inventory []Items
type Items struct {
	Type     string
	Quantity int
}

func (u User) Has(amount int) bool {
	return u.Coins >= amount
}

type WalletOperation struct {
	ID       uint
	Sender   User
	Value    int
	Receiver User
}
