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

	"golang.org/x/oauth2"
)

const (
	AnilistRedirectURL = "https://anilist.co/api/v2/oauth/pin"
	AnilistAuthURL     = "https://anilist.co/api/v2/oauth/authorize"
	AnilistTokenURL    = "https://anilist.co/api/v2/oauth/token"

	Endpoint = "https://graphql.anilist.co"
)

const (
	AnsiGreen = "\033[32m" // ANSI escape code for AnsiBlue
	AnsiBlue  = "\033[34m" // ANSI escape code for blue
	AnsiReset = "\033[0m"  // Reset color

)

//go:embed media-collection.graphql
var MediaCollectionQuery string

//go:embed user.graphql
var UserQuery string

type GraphQL struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

func (q GraphQL) Json() []byte {
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

	clientPath := filepath.Join(config, "anilist-gird", "client.json")

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

	clientPath := filepath.Join(config, "anilist-gird", "client.json")

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

func NewAnilist(ctx context.Context) *Anilist {
	cred, err := loadCredentials()
	if err != nil || cred.ID == "" || cred.Secret == "" {
		var id string
		var secret string

		fmt.Printf("%sYou need create an Anilist app!%s\n", AnsiGreen, AnsiReset)
		fmt.Printf(" - Goto %s%s%s\n", AnsiBlue, "https://anilist.co/settings/developer", AnsiReset)
		fmt.Println(" - Create New Client")
		fmt.Printf(" - Give it name and set Redirect URL to %s%s%s\n\n", AnsiBlue, AnilistRedirectURL, AnsiReset)

		fmt.Print("Enter Client ID: ")
		fmt.Scanln(&id)
		fmt.Print("Enter Client Secret: ")
		fmt.Scanln(&secret)

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

	fmt.Printf("%sOpen the following URL in your browser and authorize the application:%s\n", AnsiGreen, AnsiReset)
	fmt.Printf(" - %s%s%s\n\n", AnsiBlue, a.LoginURL(), AnsiReset)

	var code string
	fmt.Print("Paste the code: ")
	fmt.Scanln(&code)

	code = strings.TrimSpace(code)
	if code == "" {
		panic("")
	}

	if err := a.Exchange(code); err != nil {
		return err
	}

	return nil
}

func (a *Anilist) GetCurrentUser() (User, error) {
	slog.Info("Anilist.GetCurrentUser: Fetching current user")
	var user User

	query := GraphQL{Query: UserQuery, Variables: make(map[string]any)}
	jsonBytes := query.Json()

	req, err := http.NewRequestWithContext(a.ctx, http.MethodPost, Endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return user, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.http.Do(req)
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

func (a *Anilist) GetList(id int64) (AnimeList, MangaList, error) {
	type animeResult struct {
		list AnimeList
		err  error
	}

	type mangaResult struct {
		list MangaList
		err  error
	}

	animeCh := make(chan animeResult, 1)
	mangaCh := make(chan mangaResult, 1)

	// Fetch anime list concurrently
	go func() {
		defer close(animeCh)

		slog.Info("Anilist.GetList: Fetching anime list")
		animeQuery := GraphQL{Query: MediaCollectionQuery, Variables: map[string]any{"userId": id, "type": "ANIME"}}
		animeJsonBytes := animeQuery.Json()

		animeReq, err := http.NewRequestWithContext(a.ctx, http.MethodPost, Endpoint, bytes.NewBuffer(animeJsonBytes))
		if err != nil {
			animeCh <- animeResult{err: err}
			return
		}

		animeReq.Header.Set("Content-Type", "application/json")
		animeReq.Header.Set("Accept", "application/json")

		animeResp, err := a.http.Do(animeReq)
		if err != nil {
			animeCh <- animeResult{err: err}
			return
		}
		defer animeResp.Body.Close()

		if animeResp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(animeResp.Body)
			animeCh <- animeResult{err: fmt.Errorf("unexpected status code: %d, body: %s", animeResp.StatusCode, string(b))}
			return
		}

		var anime AnimeList
		if err := json.NewDecoder(animeResp.Body).Decode(&anime); err != nil {
			animeCh <- animeResult{err: err}
			return
		}

		animeCh <- animeResult{list: anime}
	}()

	// Fetch manga list concurrently
	go func() {
		defer close(mangaCh)

		slog.Info("Anilist.GetList: Fetching manga list")
		mangaQuery := GraphQL{Query: MediaCollectionQuery, Variables: map[string]any{"userId": id, "type": "MANGA"}}
		mangaJsonBytes := mangaQuery.Json()

		mangaReq, err := http.NewRequestWithContext(a.ctx, http.MethodPost, Endpoint, bytes.NewBuffer(mangaJsonBytes))
		if err != nil {
			mangaCh <- mangaResult{err: err}
			return
		}

		mangaReq.Header.Set("Content-Type", "application/json")
		mangaReq.Header.Set("Accept", "application/json")

		mangaResp, err := a.http.Do(mangaReq)
		if err != nil {
			mangaCh <- mangaResult{err: err}
			return
		}
		defer mangaResp.Body.Close()

		if mangaResp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(mangaResp.Body)
			mangaCh <- mangaResult{err: fmt.Errorf("unexpected status code: %d, body: %s", mangaResp.StatusCode, string(b))}
			return
		}

		var manga MangaList
		if err := json.NewDecoder(mangaResp.Body).Decode(&manga); err != nil {
			mangaCh <- mangaResult{err: err}
			return
		}

		mangaCh <- mangaResult{list: manga}
	}()

	// Wait for both results
	animeRes := <-animeCh
	mangaRes := <-mangaCh

	// Handle errors - return the first error encountered
	if animeRes.err != nil {
		return AnimeList{}, MangaList{}, animeRes.err
	}
	if mangaRes.err != nil {
		return AnimeList{}, MangaList{}, mangaRes.err
	}

	return animeRes.list, mangaRes.list, nil
}
