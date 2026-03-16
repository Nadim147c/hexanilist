export interface AnilistData {
  user: User
  anime: List[]
  manga: List[]
}

export interface UserResponse {
  data: { User: User }
}

export interface User {
  avatar: Avatar
  bannerImage: string
  favourites: Favourites
  id: number
  name: string
}

export interface Avatar {
  large: string
  medium: string
}

export interface Favourites {
  anime: FavouriteNode
  manga: FavouriteNode
  characters: Characters
}

export interface FavouriteNode {
  nodes: Node[]
}

export interface Node {
  id: number
}

export interface Characters {
  nodes: CharactersNode[]
}

export interface CharactersNode {
  id: number
  image: Avatar
}

export type Status = "CURRENT" | "COMPLETED" | "DROPPED" | "PAUSED" | "PLANNING"

export type MediaType = "ANIME" | "MANGA"
export interface ListData {
  MediaListCollection: MediaListCollection
}

export interface MediaListCollection {
  lists: List[]
}

export interface List {
  entries: Entry[]
  name: string
  status: Status
}

export interface Entry {
  media: Media
  score: number | null
  status: Status
}

export interface Media {
  id: number
  averageScore?: number | null
  bannerImage?: string | null
  coverImage: CoverImage
  isAdult: boolean
  meanScore: number | null
  popularity: number
  type: MediaType
}

export interface CoverImage {
  color: string | null
  extraLarge: string
  large: string
  medium: string
}

export interface AnimeList {
  data: ListData
}

export interface MangaList {
  data: ListData
}
