package service

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/gmail/v1"
)

func Init(*http.Client) {

}

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
