package utils

import (
	"fmt"

	"github.com/niluwats/bethel_dashboard/errs"
	"github.com/niluwats/bethel_dashboard/logger"

	"net/smtp"
)

func SendEmail(code string, email, name string) *errs.AppError {
	from := "testemsender@gmail.com"
	password := "ozkutuftwtsfmlxa"
	toEmail := email
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := ""
	body := ""
	if name == "userreg" {
		subject = "Subject:Email verification to signup\n\n\n"
		body = fmt.Sprintf("http://13.76.156.32:8000/auth/users/verifyemail/%s/%s", email, code)
	} else {
		subject = "Subject:Password reset link\n\n\n"
		body = fmt.Sprintf("http://13.76.156.32:8000/auth/resetpassword/%s/%s", Encode(email), code)
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
