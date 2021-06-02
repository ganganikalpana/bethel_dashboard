package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/niluwats/bethel_dashboard/errs"
	"github.com/niluwats/bethel_dashboard/logger"
)

func SendSms(to string, code int) *errs.AppError {
	to = "+94" + to
	fmt.Println("to ", to)
	secret := "b8f2241591f7552ed429e9049ae38eb6"
	key := "AC5512873a0d6528495defd49187949d08"
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", key)
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", "+13128746692")
	strCode := strconv.Itoa(code)
	msg := "your verification code is " + strCode
	fmt.Println(msg)
	msgData.Set("Body", msg)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(key, secret)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err0 := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			logger.Error(err.Error())
			return errs.NewUnexpectedError("error while decoding")
		} else {
			fmt.Println(resp.Status)
		}
	}
	if err0 != nil {
		logger.Error(err0.Error())
		return errs.NewUnexpectedError("error while sending sms---")
	} else {
		return nil
	}
}
