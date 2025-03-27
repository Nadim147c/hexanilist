# Hexanilist: Hexagon Grid Generator for Anilist

**hexanilist** is a Go-powered CLI that generates a **hexagonal grid** from your AniList profile.

<img src="./assets/hexagon.png" width="400" height="400">

## Features:

- **Your avatar** at the center.
- **Favorite characters** surrounding it.
- **High-scored completed anime/manga** forming the core.
- **Dropped and low-rated series** on the edges.

## Installation

```sh
go install github.com/Nadim147c/hexanilist@latest
```

## Usage

```sh
hexanilist -c <hexagon_size> -s <main_image_size>
```

### Options:

> Note: It would be hard to auto determine size of full grid in advanced. Therefore,
> you have to do a guessing game to get the perfect size.

- `-c int` — **Each hexagon size** (default: 50px).
- `-s int` — **Final image size** (default: 2000px).

## Example:

```sh
hexanilist -c 80 -s 2500
```

This creates a **2500px-wide** hexagon grid, with each hexagon **80px** in size.

## Why?

Hexagons offer a structured and aesthetic way to visualize your anime and manga preferences.

---

Generate your **personalized hexagonal AniList visualization** today!
