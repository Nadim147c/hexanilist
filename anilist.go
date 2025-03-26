package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/oauth2"
)

const (
	Port         = ":90123"
	CallbackPath = "/callback"

	AnilistRedirectURI = "http://localhost" + Port + CallbackPath

	AnilistAuthURL  = "https://anilist.co/api/v2/oauth/authorize"
	AnilistTokenURL = "https://anilist.co/api/v2/oauth/token"
)

// Anilist holds the OAuth2 configuration and client
type Anilist struct {
	ctx    context.Context
	oauth2 *oauth2.Config
	tok    *oauth2.Token
	http   *http.Client
}

type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func saveCredentials(clientID, clientSecret string) error {
	config, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	clientPath := filepath.Join(config, "anilist-gird", "client.json")

	dir := filepath.Dir(clientPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	credentials := Credentials{ClientID: clientID, ClientSecret: clientSecret}
	file, err := os.Create(clientPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(credentials)
}

func loadCredentials() (string, string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	clientPath := filepath.Join(config, "anilist-gird", "client.json")

	file, err := os.Open(clientPath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var credentials Credentials
	if err := json.NewDecoder(file).Decode(&credentials); err != nil {
		return "", "", err
	}

	return credentials.ClientID, credentials.ClientSecret, nil
}

func NewAnilist(ctx context.Context) *Anilist {
	clientID, clientSecret, err := loadCredentials()
	if err != nil || clientID == "" || clientSecret == "" {
		fmt.Print("Enter Client ID: ")
		fmt.Scanln(&clientID)
		fmt.Print("Enter Client Secret: ")
		fmt.Scanln(&clientSecret)

		clientID = strings.TrimSpace(clientID)
		clientSecret = strings.TrimSpace(clientSecret)

		if err := saveCredentials(clientID, clientSecret); err != nil {
			slog.Error("Failed to save credentials: %v", "error", err)
		}
	}

	oauth2 := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     oauth2.Endpoint{AuthURL: AnilistAuthURL, TokenURL: AnilistTokenURL},
		RedirectURL:  AnilistRedirectURI,
		Scopes:       []string{},
	}

	return &Anilist{ctx: ctx, oauth2: oauth2}
}

// AuthURL returns the authentication URL to redirect the user
func (a *Anilist) LoginURL() string {
	return a.oauth2.AuthCodeURL("")
}

// ExchangeCode exchanges an authorization code for an access token
func (a *Anilist) Exchange(code string) error {
	token, err := a.oauth2.Exchange(a.ctx, code)
	if err != nil {
		return err
	}
	a.tok = token

	src := a.oauth2.TokenSource(a.ctx, token)
	a.http = oauth2.NewClient(a.ctx, src)

	return err
}

func (a *Anilist) SaveToken() error {
	if a.tok == nil {
		return errors.New("Token is nil")
	}

	config, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(config, "anilist-gird", "access.json")

	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	file, err := os.Create(tokenPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(a.tok)
}

func (a *Anilist) LoadToken() (*oauth2.Token, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	tokenPath := filepath.Join(config, "anilist-gird", "access.json")

	file, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" || token.RefreshToken == "" {
		return nil, errors.New("Missing token fields")
	}

	a.tok = &token
	return &token, nil
}

func (a *Anilist) Login() error {
	if tok, err := a.LoadToken(); err == nil {
		src := a.oauth2.TokenSource(a.ctx, tok)
		a.http = oauth2.NewClient(a.ctx, src)
		return nil
	} else {
		slog.Warn("Anilist.Login: Failed to load access-token from disk", "reason", err)
	}

	codeChan := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	http.HandleFunc(CallbackPath, func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "Authorization successful, you can close this tab.")
		codeChan <- code
	})

	server := &http.Server{Addr: Port}
	go func() {
		server.ListenAndServe()
	}()

	fmt.Println("Open the following URL in your browser and authorize the application:")
	fmt.Println(a.LoginURL())

	code := <-codeChan
	srvErr := server.Shutdown(context.Background())
	if srvErr != nil {
		slog.Error("Anilist.Login: Error shutting down server", "error", srvErr)
	}

	if err := a.Exchange(code); err != nil {
		return err
	}

	return nil
}
