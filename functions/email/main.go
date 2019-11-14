package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	auth := smtp.PlainAuth(
		"",
		"apikey",
		"SG.4CGapyV1R7XXXXX",
		"smtp.sendgrid.net",
	)
	revision := "test111"
	body := fmt.Sprintf("Subject: Notification\r\n\r\nYour service %s has been created.", revision)
	err := smtp.SendMail(
		"smtp.sendgrid.net:587",
		auth,
		"sender@example.org",
		[]string{"guoyingc@cn.ibm.com"},
		[]byte(body),
	)
	if err != nil {
		log.Fatal(err)
	}
}
