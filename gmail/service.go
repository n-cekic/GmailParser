package gmail

import (
	"fmt"
	"log"

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
	return google.ConfigFromJSON(jsonKey, gmail.GmailReadonlyScope, gmail.GmailSendScope)
}
