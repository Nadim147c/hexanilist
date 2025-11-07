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
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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

var (
	CellSize = 50
	Size     = 2000
	Username = ""
	Output   = "hexagon.png"
)

func init() {
	slog.SetDefault(slog.New(log.New(os.Stderr)))

	pflag.IntVarP(&CellSize, "cell", "c", CellSize, "Size of each hexagon")
	pflag.IntVarP(&Size, "size", "s", Size, "Size of main image")
	pflag.StringVarP(&Username, "user", "u", Username, "Username of Anilist")
	pflag.StringVarP(&Output, "out", "o", Output, "Output file name")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprint(os.Stderr, "Generate Hexagon grid from anilist media\n\n")
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, pflag.CommandLine.FlagUsages())
	}

	pflag.Parse()
}

func main() {
	anilist := NewAnilist(context.Background())
	defer anilist.SaveToken()

	if err := anilist.Login(); err != nil {
		panic(err)
	}

	var user User

	if Username != "" {
		slog.Info("Fetching user data", "username", Username)
		u, err := anilist.GetUser(Username)
		if err != nil {
			panic(err)
		}
		user = u.Data.User
	} else {
		u, err := anilist.GetCurrentUser()
		if err != nil {
			panic(err)
		}
		user = u.Data.User
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

	hexs := GenerateHexagonRing(len(nodes), float64(Size/2), float64(Size/2), float64(CellSize))
	ctx := gg.NewContext(Size, Size)
	ctx.SetLineWidth(5)
	ctx.SetStrokeStyle(gg.NewSolidPattern(color.Black))

	renderHexagons(ctx, hexs, nodes)

	if !strings.HasSuffix(Output, ".png") {
		Output += ".png"
	}

	ctx.SavePNG(Output)
	slog.Info("Saving output", "output", Output, "took", time.Since(start))
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

type model struct {
	total     int
	completed int
	progress  progress.Model
	done      bool
}

type progressMsg struct{}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case progressMsg:
		m.completed++
		if m.completed >= m.total {
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.done {
		return "✅ Rendering complete!\n"
	}
	pct := float64(m.completed) / float64(m.total)
	return fmt.Sprintf(
		"Rendering hexagons... (%d/%d)\n%s\n",
		m.completed, m.total,
		m.progress.ViewAs(pct),
	)
}

func renderHexagons(ctx *gg.Context, hexs []Hexagon, nodes []HexagonNode) {
	var mu sync.Mutex

	maxConcurrent := runtime.NumCPU()
	sem := make(chan struct{}, maxConcurrent)

	p := tea.NewProgram(model{
		total:    len(hexs),
		progress: progress.New(progress.WithDefaultGradient()),
	})

	go func() {
		for i, hex := range hexs {
			sem <- struct{}{} // acquire
			go func() {
				defer func() { <-sem }() // release

				if err := renderSingleHexagon(ctx, hex, nodes[i], &mu); err != nil {
					slog.Error("Failed to render hexagon", "index", i, "error", err)
				}

				p.Send(progressMsg{})
			}()

		}
	}()

	if _, err := p.Run(); err != nil {
		slog.Error("TUI failed", "error", err)
	}
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
