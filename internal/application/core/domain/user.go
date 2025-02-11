package domain

type User struct {
	ID           uint
	Name         string
	Login        string
	PasswordHash string
	Coins        int
	Inventory    Inventory
}
type Inventory []struct {
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
