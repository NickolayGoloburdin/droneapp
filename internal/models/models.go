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
	Account struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"account"`
}
