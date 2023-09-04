package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	if !tok.Valid() {
		return nil, errors.New("oauth2 is not valid")
	}
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ch chan *oauth2.Token, config *oauth2.Config) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	// var authCode string
	// if _, err := fmt.Scan(&authCode); err != nil {
	// 	log.Fatalf("Unable to read authorization code: %v", err)
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("AAAAAAAAAAAAAAAAAAAAAAAAA")
		// Parse the authorization code from the query parameters.
		code := r.URL.Query().Get("code")

		// Exchange the authorization code for an access token.
		tok, err := config.Exchange(r.Context(), code)
		if err != nil {
			log.Fatalf("Failed to exchange code for token: %v", err)
		}
		ch <- tok
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// tok, err := config.Exchange(context.TODO(), authCode)
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve token from web: %v", err)
	// }
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatal("Unable to encode token to file")
	}
	log.Print("Token saved to file")
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config, tokenFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		log.Print(fmt.Errorf("reading token from a file failed: %w", err))
		tokenChan := make(chan *oauth2.Token, 1)
		go getTokenFromWeb(tokenChan, config)
		tok := <-tokenChan
		saveToken(tokenFile, tok)
		time.Sleep(5 * time.Second)
	}
	return config.Client(context.Background(), tok)
}
