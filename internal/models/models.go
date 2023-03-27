package models

type LoginResponse struct {
	User struct {
		Token string `json:"token"`
	} `json:"user"`
}
type SignupResponse struct {
	User struct {
		Token string `json:"token"`
	} `json:"user"`
}
type LoginRequest struct {
	Account struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"account"`
}
type SignupRequest struct {
	Account Account `json:"account"`
}
type Account struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"first_name"`
	Surname  string `json:"last_name"`
}
type TokenData struct {
	Name    string `json:"first_name"`
	Surname string `json:"last_name"`
	Email   string `json:"email"`
}
