import axios from 'axios';
import { errorValue, type ErrorValuePromize } from './error-value';
import type { ListData, UserResponse, AnilistData } from '../types/anilist';

const MEDIA_LIST_COLLECTION_QUERY = `
query MediaListCollection($type: MediaType, $userId: Int) {
  MediaListCollection(type: $type, userId: $userId) {
    lists {
      entries {
        score
        status
        media {
          id
          coverImage {
            extraLarge
            large
            medium
            color
          }
          isAdult
          type
          averageScore
          bannerImage
        }
      }
      name
      status
    }
  }
}
`;

const USER_QUERY = `
query Query($name: String) {
  User(name: $name) {
    id
    name
    avatar {
      large
      medium
    }
    bannerImage
    favourites {
      anime {
        nodes {
          id
        }
      }
      manga {
        nodes {
          id
        }
      }
      characters {
        nodes {
          id
          image {
            medium
            large
          }
        }
      }
    }
  }
}
`;

export default async function fetchAnilist(
  username: string
): ErrorValuePromize<AnilistData> {
  const [user, err] = await fetchAnilistUser(username);
  if (err !== null) {
    return [null, err];
  }
  const id = user.data.User.id;

  const [anime, animeErr] = await fetchAnilistMediaListCollection(id, 'ANIME');
  if (animeErr !== null) {
    return [null, animeErr];
  }

  const [manga, mangaErr] = await fetchAnilistMediaListCollection(id, 'MANGA');
  if (mangaErr !== null) {
    return [null, mangaErr];
  }

  const resp = {
    user: user.data.User,
    anime: anime.MediaListCollection.lists,
    manga: manga.MediaListCollection.lists,
  };

  return [resp, null];
}

type MediaType = 'ANIME' | 'MANGA';

async function fetchAnilistMediaListCollection(
  userID: number,
  mediaType: MediaType
): ErrorValuePromize<ListData> {
  const [resp, err] = await makeAnilistRequest(MEDIA_LIST_COLLECTION_QUERY, {
    userId: userID,
    type: mediaType,
  });
  if (err !== null) {
    return [null, err];
  }
  return [resp.data.data satisfies ListData, null];
}

async function fetchAnilistUser(
  username: string
): ErrorValuePromize<UserResponse> {
  const [resp, err] = await makeAnilistRequest(USER_QUERY, { name: username });
  if (err !== null) {
    return [null, err];
  }
  return [resp.data satisfies UserResponse, null];
}

async function makeAnilistRequest(
  query: string,
  variables: Record<string, any>
) {
  const request = axios.post('https://graphql.anilist.co', {
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      'Access-Control-Allow-Origin': '*',
    },
    body: JSON.stringify({ query, variables }),
  });

  return errorValue(request);
}
