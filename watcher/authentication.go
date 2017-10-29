package watcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

func getClient(secret, token string) (*http.Client, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(secret)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	t, err := tokenFromFile(token)
	if err != nil {
		t, err = tokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(token, t)
	}

	return config.Client(ctx, t), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	t, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func saveToken(file string, tok *oauth2.Token) error {
	fmt.Printf("Saving credential file to %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}
