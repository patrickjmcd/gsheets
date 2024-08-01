package gsheets

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"net/http"
	"os"
)

type ServiceAccount struct {
	Email      string `json:"client_email"`
	PrivateKey string `json:"private_key"`
}

func TryReadToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(ctx, config)
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve token from file or web")
			return nil, fmt.Errorf("unable to retrieve token from file or web: %v", err)
		}
		err = saveToken(tokFile, tok)
		if err != nil {
			log.Error().Err(err).Msg("unable to cache oauth token")
			return nil, fmt.Errorf("unable to cache oauth token: %v", err)
		}
	}
	return tok, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config, tok *oauth2.Token) (*http.Client, error) {
	return config.Client(ctx, tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Error().Err(err).Msg("unable to read authorization code")
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Error().Err(err).Msg("unable to retrieve token from web")
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func tokenFromB64(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Error().Err(err).Msg("unable to cache oauth token")
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func GetClientForServiceAccount(ctx context.Context, account ServiceAccount, scopes []string) (*http.Client, error) {
	config := &jwt.Config{
		Email:      account.Email,
		PrivateKey: []byte(account.PrivateKey),
		Scopes:     scopes,
		TokenURL:   google.JWTTokenURL,
	}
	client := config.Client(ctx)
	return client, nil
}
