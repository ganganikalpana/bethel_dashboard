package utils

import (
	"fmt"

	"github.com/niluwats/bethel_dashboard/errs"
	"github.com/niluwats/bethel_dashboard/logger"

	"net/smtp"
)

func SendEmail(code string, email, name string) *errs.AppError {
	from := "kniluwathsala@gmail.com"
	password := "rcaoznpcmmfcnqrp"
	toEmail := email
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := ""
	body := ""
	if name != "" {
		subject = "Subject:Email verification to signup\n\n\n"
		body = name + ", Your verification code is " + code
	} else {
		hashed_email := Encode(email)
		subject = "Subject:Password reset link\n\n\n"
		body = fmt.Sprintf("http://localhost:8000/auth/resetpassword/%s/%s", hashed_email, code)
	}
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, to, message)

	if err != nil {
		logger.Error(err.Error())
		return errs.NewUnexpectedError("error while sending email")
	}
	return nil
}
