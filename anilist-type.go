package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Image is a url to image
type Image string

// Download downloads the image and saves it to the disk.
// Returns the path in which the image is downloaded.
func (i Image) Download() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %w", err)
	}

	cacheDir = filepath.Join(cacheDir, "anilist-grid", "images")
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	filename := path.Base(string(i))
	filePath := filepath.Join(cacheDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		slog.Debug("Image already exists", "path", filePath)
		return filePath, nil
	}

	resp, err := http.Get(string(i))
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	slog.Debug("Image downloaded successfully", "path", filePath)
	return filePath, nil
}

// User represents the root structure of a user profile.
type User struct {
	UserData `json:"data"`
}

// UserData encapsulates the main data for a user.
type UserData struct {
	Viewer `json:"Viewer"`
}

// Viewer contains detailed information about a user, including avatar, banner, favorites, and identity.
type Viewer struct {
	Avatar     Avatar     `json:"avatar"`      // Avatar images in different sizes.
	Banner     Image      `json:"bannerImage"` // Banner image associated with the user profile.
	Favourites Favourites `json:"favourites"`  // Favorite anime, manga, and characters.
	ID         int64      `json:"id"`          // Unique identifier of the user.
	Name       string     `json:"name"`        // Display name of the user.
}

// Avatar holds different sizes of an avatar image.
type Avatar struct {
	Large  Image `json:"large"`  // Large-sized avatar image.
	Medium Image `json:"medium"` // Medium-sized avatar image.
}

// Favourites stores the user's favorite anime, manga, and characters.
type Favourites struct {
	Anime      FavouriteNode `json:"anime"`      // Favorite anime list.
	Manga      FavouriteNode `json:"manga"`      // Favorite manga list.
	Characters Characters    `json:"characters"` // Favorite characters.
}

// FavouriteNode represents a list of favorite anime or manga.
type FavouriteNode struct {
	Nodes `json:"nodes"` // Collection of favorite media entries.
}

func (fn FavouriteNode) Has(id int64) bool {
	for _, n := range fn.Nodes {
		if n.ID == id {
			return true
		}
	}
	return false
}

// Nodes is a slice of Node, representing a collection of media entries.
type Nodes []Node

// Node represents a single media entry.
type Node struct {
	ID int64 `json:"id"` // Unique identifier of the media entry.
}

// Characters represents a collection of favorite character nodes.
type Characters struct {
	Nodes []CharactersNode `json:"nodes"` // List of favorite character entries.
}

// CharactersNode represents a single favorite character entry.
type CharactersNode struct {
	ID    int64  `json:"id"`    // Unique identifier of the character.
	Image Avatar `json:"image"` // Character's avatar image.
}

// Status represents different statuses for a media list entry.
type Status string

const (
	Current   Status = "CURRENT"   // Currently watching/reading.
	Completed Status = "COMPLETED" // Fully watched/read.
	Dropped   Status = "DROPPED"   // Dropped midway.
	Paused    Status = "PAUSED"    // Temporarily on hold.
	Planning  Status = "PLANNING"  // Planned for future.
)

// Type defines whether the media is anime or manga.
type Type string

const (
	Anime Type = "ANIME" // Anime type media.
	Manga Type = "MANGA" // Manga type media.
)

// ListData represents the root structure for media lists.
type ListData struct {
	MediaListCollection `json:"MediaListCollection"`
}

// MediaListCollection contains a list of categorized media lists.
type MediaListCollection struct {
	Lists []List `json:"lists"` // Collection of media lists.
}

// List represents a categorized list of media entries.
type List struct {
	Entries []Entry `json:"entries"` // Entries in the list.
	Name    string  `json:"name"`    // Name of the list.
	Status  Status  `json:"status"`  // Status of the list.
}

// Entry represents a single media entry with a score.
type Entry struct {
	Media `json:"media"` // Media details.

	Score  *float64 `json:"score"`  // User-assigned score.
	Status Status   `json:"status"` // User-assigned status.
}

// Media represents detailed information about a media entry.
type Media struct {
	ID           int64      `json:"id"`           //  ID of the media
	AverageScore *int64     `json:"averageScore"` // Community average score.
	Banner       *Image     `json:"bannerImage"`  // Banner image of the media.
	Cover        CoverImage `json:"coverImage"`   // Cover image in different sizes.
	IsAdult      bool       `json:"isAdult"`      // Indicates if the media is for adults.
	MeanScore    *int64     `json:"meanScore"`    // Mean score of the media.
	Popularity   int64      `json:"popularity"`   // Popularity ranking.
	Type         Type       `json:"type"`         // Media type (anime/manga).
}

// CoverImage contains multiple sizes of the cover image.
type CoverImage struct {
	Color      *string `json:"color"`      // Dominant color of the cover image.
	ExtraLarge Image   `json:"extraLarge"` // Extra-large-sized cover image.
	Large      Image   `json:"large"`      // Large-sized cover image.
	Medium     Image   `json:"medium"`     // Medium-sized cover image.
}

// AnimeList is a wrapper for anime-related media lists.
type AnimeList struct {
	ListData `json:"data"`
}

// MangaList is a wrapper for manga-related media lists.
type MangaList struct {
	ListData `json:"data"`
}
