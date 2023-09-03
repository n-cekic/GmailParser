package gmail

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type Service struct {
	Srv *gmail.Service
}

// RetrieveLabels retrieve labels from gmail service
func (s *Service) RetrieveLabels() {
	user := "me"
	response, err := s.Srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	if len(response.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}

	fmt.Println("Labels:")
	for _, l := range response.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}

func GetConfigFromJSON(jsonKey []byte, scope ...string) (*oauth2.Config, error) {
	return google.ConfigFromJSON(jsonKey, gmail.GmailReadonlyScope, gmail.GmailSendScope, gmail.GmailModifyScope)
}
func CheckForNewMessages(srv *gmail.Service) {
	previouslyPrinted := false
	for {
		messages, err := retrieveUnreadMessages(*srv)
		if err != nil {
			log.Print(fmt.Errorf("failed listing unread messages: %w", err))
		}

		if len(messages) == 0 && !previouslyPrinted {
			log.Print("No new messages")
			previouslyPrinted = true
		}

		if len(messages) > 0 {
			previouslyPrinted = false
		}

		for _, msg := range messages {
			message, err := getFullMessage(srv, msg.Id)
			if err != nil {
				log.Print(fmt.Errorf("failed getting full message: %w", err))
			}
			for _, hdr := range message.Payload.Headers {
				if hdr.Name == "Subject" {
					log.Print("subject of this message is: ", hdr.Value)
				}
			}
			data, err := decodeMessagePart(message.Payload.Parts[0])
			if err != nil {
				log.Print(fmt.Errorf("failed decoding message payload: %w", err))
			}

			log.Print(data)

			markMessageAsRead(srv, message)
		}

		time.Sleep(5 * time.Second)
	}
}

func retrieveUnreadMessages(srv gmail.Service) ([]*gmail.Message, error) {
	messages, err := srv.Users.Messages.List("me").Q("is:unread").Do()
	if err != nil {
		return nil, err
	}
	return messages.Messages, nil
}

func getFullMessage(srv *gmail.Service, messageID string) (*gmail.Message, error) {
	message, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func decodeMessagePart(part *gmail.MessagePart) (string, error) {
	data, err := base64.URLEncoding.DecodeString(part.Body.Data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func markMessageAsRead(srv *gmail.Service, message *gmail.Message) {
	modifyRequest := gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{"UNREAD"},
	}

	messageID := message.Id
	_, err := srv.Users.Messages.Modify("me", messageID, &modifyRequest).Do()
	if err != nil {
		log.Println(fmt.Errorf("failed marking message %s as read: %w", messageID, err))
	}
	log.Printf("Message %s marked as read", messageID)
}
