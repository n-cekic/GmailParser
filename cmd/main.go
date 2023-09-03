package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"nikola/monygo/auth"
	mailing "nikola/monygo/gmail"
)

var secretsPath = flag.String("token.file", "${fileDirname}/../secrets/", "path to token file")

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
	flag.Parse()

	b, err := os.ReadFile(*secretsPath + "client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := mailing.GetConfigFromJSON(b)
	if err != nil {
		log.Fatal(fmt.Errorf("failed creating client configuration: %w", err))
	}

	client := auth.GetClient(config, *secretsPath+"token.json")

	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	go mailing.CheckForNewMessages(srv)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	<-quitChannel
}
