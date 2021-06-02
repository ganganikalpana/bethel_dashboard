package dto

type NewLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
