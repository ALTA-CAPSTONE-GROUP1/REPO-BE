package helper

import (
	"fmt"
	"net/smtp"
)

func SendSimpleWhatsappNotif(phoneNumber, registrationData string) error {
	// Replace the placeholders below with your email server configuration
	smtpHost := "your-smtp-host"
	smtpPort := 587
	senderEmail := "sender@example.com"
	senderPassword := "password"

	// Compose the email message
	message := fmt.Sprintf("Subject: Registration Details\n\nThank you for registering. Here are your registration details:\n%s", registrationData)

	// Compose the recipient's email address using the phone number
	recipientEmail := phoneNumber + "@provider.com"

	// Configure the SMTP client
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Send the email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, senderEmail, []string{recipientEmail}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
