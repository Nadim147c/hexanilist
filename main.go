package main

import (
	"context"
	"fmt"
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

	for _, list := range anime.Lists {
		fmt.Printf("%s := %+v\n", list.Name, len(list.Entries))
	}
	for _, list := range manga.Lists {
		fmt.Printf("%s := %+v\n", list.Name, len(list.Entries))
	}
}
