package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/oauth2"
)

const (
	AnilistDeveloperPortal = "https://anilist.co/settings/developer"
	AnilistRedirectURL     = "https://anilist.co/api/v2/oauth/pin"
	AnilistAuthURL         = "https://anilist.co/api/v2/oauth/authorize"
	AnilistTokenURL        = "https://anilist.co/api/v2/oauth/token"

	Endpoint = "https://graphql.anilist.co"
)

//go:embed media-collection.graphql
var MediaCollectionQuery string

//go:embed viewer.graphql
var ViewerQuery string

//go:embed user.graphql
var UserQuery string

type GraphQL struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func (q GraphQL) JSON() []byte {
	b, err := json.Marshal(q)
	if err != nil {
		slog.Error("Query.String: Failed to Marshal query")
		return []byte{}
	}
	return b
}

// Anilist holds the OAuth2 configuration and client
type Anilist struct {
	ctx    context.Context
	oauth2 *oauth2.Config
	tok    *oauth2.Token
	http   *http.Client
}

type Credentials struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

func saveCredentials(credentials Credentials) error {
	config, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	clientPath := filepath.Join(config, "hexanilist", "client.json")

	dir := filepath.Dir(clientPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	file, err := os.Create(clientPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(credentials)
}

func loadCredentials() (Credentials, error) {
	var credentials Credentials
	config, err := os.UserConfigDir()
	if err != nil {
		return credentials, err
	}

	clientPath := filepath.Join(config, "hexanilist", "client.json")

	file, err := os.Open(clientPath)
	if err != nil {
		return credentials, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&credentials); err != nil {
		return credentials, err
	}

	return credentials, nil
}

var (
	TitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)   // green color
	URLStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Italic(true) // blue color
)

func NewAnilist(ctx context.Context) *Anilist {
	cred, err := loadCredentials()
	if err != nil || cred.ID == "" || cred.Secret == "" {
		var id string
		var secret string

		fmt.Printf(
			"%s\n - Goto %s\n - Create A New Client\n - Give it name and set Redirect URL to %s\n\n",
			TitleStyle.Render("You need create an Anilist app!"),
			URLStyle.Render(AnilistDeveloperPortal),
			URLStyle.Render(AnilistRedirectURL),
		)

		huh.NewInput().
			Title("Anilist API Client ID").
			Value(&id).
			Run()

		fmt.Println("Client ID: ", id)

		huh.NewInput().
			Title("Anilist API Client Secret").
			EchoMode(huh.EchoModePassword).
			Value(&secret).
			Run()

		cred = Credentials{
			ID:     strings.TrimSpace(id),
			Secret: strings.TrimSpace(secret),
		}

		if cred.ID == "" || cred.Secret == "" {
			panic("Invalid ID and Secret")
		}

		if err := saveCredentials(cred); err != nil {
			slog.Error("Failed to save credentials: %v", "error", err)
		}
	}

	oauth2 := &oauth2.Config{
		ClientID:     cred.ID,
		ClientSecret: cred.Secret,
		Endpoint:     oauth2.Endpoint{AuthURL: AnilistAuthURL, TokenURL: AnilistTokenURL},
		RedirectURL:  AnilistRedirectURL,
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

	tokenPath := filepath.Join(config, "hexanilist", "access.json")

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

	tokenPath := filepath.Join(config, "hexanilist", "access.json")

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

	fmt.Printf(
		"%s\n - %s\n\n",
		URLStyle.Render(a.LoginURL()),
		TitleStyle.Render("Open the following URL in your browser and authorize the application"),
	)

	var code string
	fmt.Print("Paste the code: ")
	huh.NewInput().
		Title("Paste the code").
		EchoMode(huh.EchoModePassword).
		Value(&code).
		Run()

	code = strings.TrimSpace(code)
	if code == "" {
		panic("")
	}

	if err := a.Exchange(code); err != nil {
		return err
	}

	return nil
}

func (a *Anilist) makeRequest(q GraphQL) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		a.ctx, http.MethodPost, Endpoint,
		bytes.NewBuffer(q.JSON()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return a.http.Do(req)
}

func (a *Anilist) GetCurrentUser() (Viewer, error) {
	slog.Info("Anilist.GetCurrentUser: Fetching current user")
	var user Viewer

	query := GraphQL{Query: ViewerQuery, Variables: make(map[string]any)}
	resp, err := a.makeRequest(query)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (a *Anilist) GetUser(username string) (Searched, error) {
	slog.Info("Anilist.GetCurrentUser: Fetching current user")
	var user Searched

	query := GraphQL{
		Query:     UserQuery,
		Variables: map[string]any{"name": username},
	}

	resp, err := a.makeRequest(query)
	if err != nil {
		return user, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (a *Anilist) GetList(id int64) (anime AnimeList, manga MangaList, err error) {
	slog.Info("Anilist.GetList: Fetching anime list")

	// --- Fetch Anime List ---
	queryAnime := GraphQL{
		Query:     MediaCollectionQuery,
		Variables: map[string]any{"userId": id, "type": "ANIME"},
	}

	resp, err := a.makeRequest(queryAnime)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return anime, manga, fmt.Errorf(
			"unexpected status code (anime): %d, body: %s",
			resp.StatusCode,
			string(b),
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&anime); err != nil {
		return anime, manga, err
	}

	// --- Fetch Manga List ---
	slog.Info("Anilist.GetList: Fetching manga list")

	queryManga := GraphQL{
		Query:     MediaCollectionQuery,
		Variables: map[string]any{"userId": id, "type": "MANGA"},
	}

	resp, err = a.makeRequest(queryManga)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return anime, manga, fmt.Errorf(
			"unexpected status code (manga): %d, body: %s",
			resp.StatusCode,
			string(b),
		)
	}

	if err := json.NewDecoder(resp.Body).Decode(&manga); err != nil {
		return anime, manga, err
	}

	return
}
