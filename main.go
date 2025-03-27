package main

import (
	"context"
	"fmt"
	"image/color"
	"log/slog"
	"slices"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type NodeType uint64

const (
	UserNode NodeType = iota
	AnimeNode
	MangaNode
	CharacterNode
)

type (
	HexagonNode struct {
		Type  NodeType
		Image Image
		Score int64
	}
)

func main() {
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

	userNode := HexagonNode{
		Type:  UserNode,
		Score: 1 << 60,
		Image: user.Avatar.Medium,
	}

	nodes := []HexagonNode{userNode}

	for _, char := range user.Favourites.Characters.Nodes {
		characterNode := HexagonNode{
			Type:  CharacterNode,
			Score: 500,
			Image: char.Image.Medium,
		}

		nodes = append(nodes, characterNode)
	}

	for _, list := range anime.Lists {
		for _, entry := range list.Entries {
			var score int64 = 0
			if entry.Score != nil {
				score = int64(*entry.Score * 10)
			}

			switch entry.Status {
			case Completed:
				score += 100
			case Dropped:
				score -= 100
			}

			if user.Favourites.Manga.Has(entry.ID) {
				score += 200
			}

			animeNode := HexagonNode{
				Type:  AnimeNode,
				Score: score,
				Image: entry.Cover.Medium,
			}

			nodes = append(nodes, animeNode)
		}
	}

	for _, list := range manga.Lists {
		for _, entry := range list.Entries {
			var score int64 = 0
			if entry.Score != nil {
				score = int64(*entry.Score * 10)
			}

			switch entry.Status {
			case Completed:
				score += 100
			case Dropped:
				score -= 100
			}

			if user.Favourites.Manga.Has(entry.ID) {
				score += 200
			}

			animeNode := HexagonNode{
				Type:  AnimeNode,
				Score: score,
				Image: entry.Cover.Medium,
			}

			nodes = append(nodes, animeNode)
		}
	}

	slices.SortFunc(nodes, func(i, j HexagonNode) int {
		if i.Score > j.Score {
			return -1
		}
		if i.Score < j.Score {
			return 1
		}
		return 0
	})

	w, h := 2000, 2000
	hexs := GenerateHexagonRing(len(nodes), float64(w/2), float64(h/2), 40)

	ctx := gg.NewContext(w, h)
	ctx.SetLineWidth(5)
	ctx.SetStrokeStyle(gg.NewSolidPattern(color.Black))

	for i, hex := range hexs {
		slog.Info("Process image", "#", i)
		node := nodes[i]
		path, err := node.Image.Download()
		if err != nil {
			slog.Error("Failed to download image.", "error", err)
			continue
		}

		img, err := gg.LoadImage(path)
		if err != nil {
			slog.Error("Failed to load image.", "error", err)
			continue
		}

		w, h := hex.Box().Size()
		cropped := imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)

		ctx.ClearPath()
		hex.Draw(ctx)
		ctx.Clip()

		x, y := hex.Box().Start()
		ctx.DrawImage(cropped, x, y)
		ctx.ResetClip()

		// Draw stroke arround the image
		hex.Draw(ctx)
		ctx.Stroke()
	}

	ctx.SavePNG("hexagon.png")
}
