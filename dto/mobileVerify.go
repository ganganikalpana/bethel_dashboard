package dto

type NewMobileVerificationRequest struct {
	Mobile string `json:"mobile"`
}
type MobileVerifyCode struct {
	Mobile string `json:"mobile"`
	Code   int    `json:"code"`
}
