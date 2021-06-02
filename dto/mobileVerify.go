package dto

type MobileVerifyCode struct {
	Mobile string `json:"mobile"`
	Code   int    `json:"code"`
}
