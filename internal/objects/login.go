package objects

type LoginMessage struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenMessage struct {
	Token string `json:"token"`
}
