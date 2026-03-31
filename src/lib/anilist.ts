import axios from "axios"
import type {
  ListData,
  UserResponse,
  AnilistData,
  Status,
} from "../types/anilist"
import { handlePromise, type AsyncResult } from "./result"
import { MEDIA_LIST_COLLECTION_QUERY, USER_QUERY } from "./anilist-queries"

interface Node {
  type: "user" | "anime" | "manga" | "character"
  image: string
  score: number
}

function calculateScore(
  userScore: number | null,
  status: Status,
  isFavorite: boolean
): number {
  let score = 0
  if (userScore !== null) score = userScore * 10
  if (isFavorite) score += 200
  switch (status) {
    case "COMPLETED":
      score += 100
      break
    case "DROPPED":
      score -= 100
      break
  }
  return score
}

export function createNodes(data: AnilistData): Node[] {
  const nodes: Node[] = []
  nodes.push({
    type: "user",
    image: data.user.avatar.medium,
    score: 10000000,
  })

  data.user.favourites.characters.nodes.forEach(character => {
    nodes.push({
      type: "character",
      image: character.image.medium,
      score: 5000,
    })
  })

  data.anime.forEach(list => {
    list.entries.forEach(entry => {
      const isFavorite = data.user.favourites.anime.nodes.some(
        node => node.id === entry.media.id
      )
      nodes.push({
        type: "anime",
        image: entry.media.coverImage.medium,
        score: calculateScore(entry.score, entry.status, isFavorite),
      })
    })
  })

  data.manga.forEach(list => {
    list.entries.forEach(entry => {
      const isFavorite = data.user.favourites.manga.nodes.some(
        node => node.id === entry.media.id
      )
      nodes.push({
        type: "manga",
        image: entry.media.coverImage.medium,
        score: calculateScore(entry.score, entry.status, isFavorite),
      })
    })
  })

  return nodes
}

export async function fetchAnilist(username: string): AsyncResult<AnilistData> {
  const [user, err] = await fetchAnilistUser(username)
  if (err !== null) {
    return [null, err]
  }
  const id = user.data.User.id

  const [anime, animeErr] = await fetchAnilistMediaListCollection(id, "ANIME")
  if (animeErr !== null) {
    return [null, animeErr]
  }

  const [manga, mangaErr] = await fetchAnilistMediaListCollection(id, "MANGA")
  if (mangaErr !== null) {
    return [null, mangaErr]
  }

  const resp = {
    user: user.data.User,
    anime: anime.MediaListCollection.lists,
    manga: manga.MediaListCollection.lists,
  }

  return [resp, null]
}

type MediaType = "ANIME" | "MANGA"

async function fetchAnilistMediaListCollection(
  userID: number,
  mediaType: MediaType
): AsyncResult<ListData> {
  const [resp, err] = await makeAnilistRequest(MEDIA_LIST_COLLECTION_QUERY, {
    userId: userID,
    type: mediaType,
  })
  if (err !== null) {
    return [
      null,
      new Error(`Failed to fetch Anilist ${mediaType}: ${err.message}`),
    ]
  }
  return [resp.data.data satisfies ListData, null]
}

async function fetchAnilistUser(username: string): AsyncResult<UserResponse> {
  const [resp, err] = await makeAnilistRequest(USER_QUERY, { name: username })
  if (err !== null) {
    return [null, new Error(`Failed to fetch Anilist User: ${err.message}`)]
  }
  return [resp.data satisfies UserResponse, null]
}

async function makeAnilistRequest(
  query: string,
  variables: Record<string, any>
) {
  const head = {
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  }
  const payload = { query, variables }
  const request = axios.post("https://graphql.anilist.co", payload, head)
  return handlePromise(request)
}
