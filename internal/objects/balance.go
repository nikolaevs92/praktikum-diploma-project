package objects

type Balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdrawn"`
}
