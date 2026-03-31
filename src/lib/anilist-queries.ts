export const MEDIA_LIST_COLLECTION_QUERY = `
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
`

export const USER_QUERY = `
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
`
