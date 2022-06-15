package objects

type RegisterMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
