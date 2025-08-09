package user

type Item struct {
	Username string `json:"Username"`
	Email    string `json:"email"`
}

type ListResponse struct {
	Users []Item `json:"users"`
	F1    string `json:"f1"`
}
