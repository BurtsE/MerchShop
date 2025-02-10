package domain

type Merch struct {
	Name string
	Cost int64
}

type User struct {
	ID           uint
	Login        string
	PasswordHash string
}

type Transfer struct {
	ID       uint
	FromUser uint
	ToUser   uint
}
