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

	fmt.Printf("%+v\n", user)
}
