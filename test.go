package main

import (
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"log"
	"net/http"
)

func main() {
	// Initialize your Gmail API configuration (config) and client (client).

	// Define your OAuth2 callback URL. This should match the URL you specified
	// when creating your OAuth2 credentials in the Google Cloud Console.

	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		// Parse the authorization code from the query parameters.
		code := r.URL.Query().Get("code")

		// Exchange the authorization code for an access token.
		token, err := config.Exchange(r.Context(), code)
		if err != nil {
			log.Fatalf("Failed to exchange code for token: %v", err)
		}

		// Now you have the access token in the 'token' variable.
		// You can use this token to make authenticated requests to the Gmail API.

		// Example: Use the token to create a Gmail API service client.
		gmailService, err := gmail.New(client)
		if err != nil {
			log.Fatalf("Failed to create Gmail service client: %v", err)
		}

		// You can now use the 'gmailService' to access the Gmail API on behalf of the user.

		// Handle the response, e.g., display a success page or redirect the user.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authentication successful! You can now access Gmail on behalf of the user."))
	})

	// Start your HTTP server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
