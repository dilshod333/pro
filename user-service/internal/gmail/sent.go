package email

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"gopkg.in/gomail.v2"
)

func sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "dilshoddilmurodov112@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dilshoddilmurodov112@gmail.com", "xmxu rdhp pmdf pezk")

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmail(to string) (code string) {
	code, err := generateRandomCode()
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	subject := "----Welcome buddy----"
	body := fmt.Sprintf("Your confirmation code is: %s", code)

	if err := sendEmail(to, subject, body); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	return code

}

func generateRandomCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())
	return code, nil
}
