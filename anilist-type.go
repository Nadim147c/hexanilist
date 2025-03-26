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

type Image string

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

type User struct {
	Data `json:"data"`
}

type Data struct {
	Viewer `json:"Viewer"`
}

type Viewer struct {
	Avatar     Avatar     `json:"avatar"`
	Banner     Image      `json:"bannerImage"`
	Favourites Favourites `json:"favourites"`
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
}

type Avatar struct {
	Large  Image `json:"large"`
	Medium Image `json:"medium"`
}

type Favourites struct {
	Anime      Anime      `json:"anime"`
	Characters Characters `json:"characters"`
	Manga      Anime      `json:"manga"`
}

type Anime struct {
	Nodes []AnimeNode `json:"nodes"`
}

type AnimeNode struct {
	ID int64 `json:"id"`
}

type Characters struct {
	Nodes []CharactersNode `json:"nodes"`
}

type CharactersNode struct {
	ID    int64  `json:"id"`
	Image Avatar `json:"image"`
}
