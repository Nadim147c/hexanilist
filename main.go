package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/spf13/pflag"
)

type NodeType uint64

const (
	UserNode NodeType = iota
	AnimeNode
	MangaNode
	CharacterNode
)

type HexagonNode struct {
	Type  NodeType
	Image Image
	Score int
}

func main() {
	cellSize := pflag.IntP("cell", "c", 50, "Size of each hexagon")
	size := pflag.IntP("size", "s", 2000, "Size of main image")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	anilist := NewAnilist(context.Background())
	defer anilist.SaveToken()

	if err := anilist.Login(); err != nil {
		panic(err)
	}

	user, err := anilist.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("User Name:", user.Name)

	anime, manga, err := anilist.GetList(user.ID)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	nodes := buildNodes(user, anime, manga)
	slices.SortFunc(nodes, func(i, j HexagonNode) int { return j.Score - i.Score })

	hexs := GenerateHexagonRing(len(nodes), float64(*size/2), float64(*size/2), float64(*cellSize))
	ctx := gg.NewContext(*size, *size)
	ctx.SetLineWidth(5)
	ctx.SetStrokeStyle(gg.NewSolidPattern(color.Black))

	renderHexagons(ctx, hexs, nodes)

	ctx.SavePNG("hexagon.png")
	slog.Info("Saved to hexagon.png", "took", time.Since(start))
}

func buildNodes(user User, anime AnimeList, manga MangaList) []HexagonNode {
	userNode := HexagonNode{
		Type:  UserNode,
		Score: 1 << 60,
		Image: user.Avatar.Medium,
	}

	nodes := []HexagonNode{userNode}
	nodes = append(nodes, buildCharacterNodes(user)...)

	nodeChan := make(chan HexagonNode)
	var wg sync.WaitGroup

	wg.Add(2)
	go processAnimeList(anime, user, nodeChan, &wg)
	go processMangaList(manga, user, nodeChan, &wg)

	go func() {
		wg.Wait()
		close(nodeChan)
	}()

	for node := range nodeChan {
		nodes = append(nodes, node)
	}

	return nodes
}

func buildCharacterNodes(user User) []HexagonNode {
	var nodes []HexagonNode
	for _, char := range user.Favourites.Characters.Nodes {
		characterNode := HexagonNode{
			Type:  CharacterNode,
			Score: 500,
			Image: char.Image.Medium,
		}
		nodes = append(nodes, characterNode)
	}
	return nodes
}

func processAnimeList(anime AnimeList, user User, nodeChan chan<- HexagonNode, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, list := range anime.Lists {
		for _, entry := range list.Entries {
			score := calculateScore(entry.Score, entry.Status, user.Favourites.Anime.Has(entry.ID))

			animeNode := HexagonNode{
				Type:  AnimeNode,
				Score: score,
				Image: entry.Cover.Medium,
			}

			nodeChan <- animeNode
		}
	}
}

func processMangaList(manga MangaList, user User, nodeChan chan<- HexagonNode, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, list := range manga.Lists {
		for _, entry := range list.Entries {
			score := calculateScore(entry.Score, entry.Status, user.Favourites.Manga.Has(entry.ID))

			mangaNode := HexagonNode{
				Type:  MangaNode, // Fixed: was AnimeNode, should be MangaNode
				Score: score,
				Image: entry.Cover.Medium,
			}

			nodeChan <- mangaNode
		}
	}
}

func calculateScore(userScore *float64, status Status, isFavorite bool) int {
	var score int = 0

	if userScore != nil {
		score = int(*userScore * 10)
	}

	switch status {
	case Completed:
		score += 100
	case Dropped:
		score -= 100
	}

	if isFavorite {
		score += 200
	}

	return score
}

func renderHexagons(ctx *gg.Context, hexs []Hexagon, nodes []HexagonNode) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	maxConcurrent := runtime.NumCPU()
	sem := make(chan struct{}, maxConcurrent)

	for i, hex := range hexs {
		wg.Add(1)
		sem <- struct{}{} // acquire

		go func(i int, hex Hexagon) {
			defer wg.Done()
			defer func() { <-sem }() // release

			if err := renderSingleHexagon(ctx, hex, nodes[i], &mu); err != nil {
				slog.Error("Failed to render hexagon", "index", i, "error", err)
			}
		}(i, hex)
	}

	wg.Wait()
}

func renderSingleHexagon(ctx *gg.Context, hex Hexagon, node HexagonNode, mu *sync.Mutex) error {
	path, err := node.Image.Download()
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	img, err := gg.LoadImage(path)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	w, h := hex.Box().Size()
	cropped := imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)

	mu.Lock()
	defer mu.Unlock()

	drawHexagonWithImage(ctx, hex, cropped)
	return nil
}

func drawHexagonWithImage(ctx *gg.Context, hex Hexagon, img image.Image) {
	// Clear path and set clipping
	ctx.ClearPath()
	hex.Draw(ctx)
	ctx.Clip()

	// Draw the image
	x, y := hex.Box().Start()
	ctx.DrawImage(img, x, y)
	ctx.ResetClip()

	// Draw stroke around the image
	hex.Draw(ctx)
	ctx.Stroke()
}
