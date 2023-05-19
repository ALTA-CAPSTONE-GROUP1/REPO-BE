package helper

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"gopkg.in/gomail.v2"
)

func SendSimpleEmail(subTitle string, subject string, recipientEmail []string, recipentName []string, senderName string) {

	mailer := gomail.NewDialer("smtp.gmail.com", 587, "eprop.unit3@gmail.com", "ozlxhlzyywzoenji")

	for i := 0; i < len(recipientEmail); i++ {
		body := "Hello! " + recipentName[i] + "\n\nThere is an update for you eProposal Account! \n\nPlease Check it Out! \n\nFrom: " + senderName

		email := gomail.NewMessage()
		email.SetHeader("From", "eprop.unit3@gmail.com")
		email.SetHeader("To", recipientEmail[i])
		email.SetHeader("Subject", subject)
		email.SetBody("text/plain", body)

		err := mailer.DialAndSend(email)
		if err != nil {
			log.Errorf("error on sending email %w", err)
		} else {
			fmt.Println("succes to send email!")
		}
	}
}
