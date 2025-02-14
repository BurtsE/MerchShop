package router

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendCoinsInfo struct {
	Username string `json:"toUser"`
	Amount   int    `json:"amount"`
}
