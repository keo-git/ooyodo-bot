package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/keo-git/ooyodo-bot/config"
	"github.com/keo-git/ooyodo-bot/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func getClient(ctx context.Context, oauthConfig *oauth2.Config) *http.Client {
	conf := config.Config()
	tokenFile, err := utils.AbsolutePath(conf.Credentials, conf.GmailToken)
	if err != nil {
		log.Fatalf("Unable to get path to token file: %v", err)
	}
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = tokenFromWeb(oauthConfig)
		saveToken(tokenFile, tok)
	}
	return oauthConfig.Client(ctx, tok)
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

func tokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return tok
}

func saveToken(file string, tok *oauth2.Token) {
	fmt.Printf("Saving credential file to %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
}
