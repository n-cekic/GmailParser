package main

import (
	"context"
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
	b, err := os.ReadFile("../credentials.json")
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

	/*
		for false {
			lbl, err := ser.Srv.Users.Labels.Get("me", "INBOX").Do()
			if err != nil {
				log.Printf("failed retrieving messages with label \"INBOX\": %s", err.Error())
			}
			log.Printf("=== unread messages from inbox: %d ===\n", lbl.MessagesUnread)

			time.Sleep(5 * time.Second)
		}
	*/
	messagesw, err := srv.Users.Messages.List("me").Q("is:unread").Do()
	if err != nil {
		log.Print(fmt.Errorf("failed listing unread messages: %w", err))
	}
	for _, messagew := range messagesw.Messages {
		message, err := srv.Users.Messages.Get("me", messagew.Id).Format("full").Do()
		if err != nil {
			log.Print(fmt.Errorf("failed getting full message: %w", err))
		}

		var content string
		for _, part := range message.Payload.Parts {
			data, err := decodeMessagePart(part)
			if err != nil {
				log.Print(fmt.Errorf("failed decoding message payload: %w", err))
			}
			content += data
		}

		fmt.Println(content)
	}

	/*
		if err != nil {
			log.Print(err)
		}
		for _, mes := range messagesList.Messages {
			log.Printf("mes.Raw: %v\n", mes.Payload.Body.Data)
		}
		srv.Users.Messages.Get("me", ";lk")
	*/

	// create an email
	// emailContent := createContent()

	// Encode the email content as base64
	// rawMessage := base64.URLEncoding.EncodeToString([]byte(emailContent))

	// msg := createEmail(rawMessage)

	// sendEmail(srv, msg)
}

func decodeMessagePart(part *gmail.MessagePart) (string, error) {
	data, err := base64.URLEncoding.DecodeString(part.Body.Data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
