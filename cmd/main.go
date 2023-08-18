package main

import (
	"context"
	//"encoding/base64"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"nikola/monygo/auth"
	mailing "nikola/monygo/gmail"
)

func sendEmail(srv *gmail.Service, msg gmail.Message) {
	// send an email
	_, err := srv.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		log.Fatalf("Unable to send message: %v", err)
	}
	fmt.Println("Email sent successfully!")
}

func createEmail(rawMessage string) gmail.Message {
	return gmail.Message{
		Raw: rawMessage,
	}
}

func createContent() string {
	from := "nikolafordev@gmail.com" // Replace with the sender's email address
	to := "ncekic13@gmail.com"       // Replace with the recipient's email address
	subject := "Test Email from Gmail API"
	message := "This is a test email sent using the Gmail API in Golang!"

	// Create the email payload
	return "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := mailing.GetConfigFromJSON(b)
	if err != nil {
		log.Fatal(fmt.Errorf("failed creating client configuration: %w", err))
	}

	client := auth.GetClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	ser := mailing.Service{Srv: srv}

	// retrieve labels
	ser.RetrieveLabels()

	// create an email
	emailContent := createContent()

	// Encode the email content as base64
	rawMessage := base64.URLEncoding.EncodeToString([]byte(emailContent))

	msg := createEmail(rawMessage)

	sendEmail(srv, msg)
}
