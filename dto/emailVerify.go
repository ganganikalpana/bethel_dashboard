package dto

type EmailVerifyCode struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}
