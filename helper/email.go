package helper

import (
	"fmt"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/labstack/gommon/log"
	"gopkg.in/gomail.v2"
)

func SendSimpleEmail(sign string, subTitle string, subject string, recipientEmail []string, recipientNames []string, senderName string) {
	mailer := gomail.NewDialer("smtp.gmail.com", 465, config.Email, config.EmailSecret)
	log.Error(config.Email)
	log.Error(config.EmailSecret)
	for i := 0; i < len(recipientEmail); i++ {
		var recipientName string
		if len(recipientNames) <= i {
			recipientName = ""
		} else {
			recipientName = recipientNames[i]
		}
		body := fmt.Sprintf(`%s
Hello, %s!

There is an update for your eProposal Account!

Please check it out.

From: %s

WITH SIGN ID : %s
`, subTitle, recipientName, senderName, sign)

		email := gomail.NewMessage()
		email.SetHeader("From", config.Email)

		email.SetHeader("To", recipientEmail[i])
		email.SetHeader("Subject", subject)
		email.SetBody("text/plain", body)

		err := mailer.DialAndSend(email)
		if err != nil {
			log.Errorf("Error sending email: %v", err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}
}
