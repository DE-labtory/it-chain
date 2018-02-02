package main

import (
	"fmt"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"context"
)

func main() {
	ctx := context.Background()
	c := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "31d8f4c1bfc6906806b9a77803087b5b671fac2d"},
	))
	client := github.NewClient(c)

	repo, res, err := client.Repositories.Create(context.Background(), "", &github.Repository{Name: github.String("test1"), AutoInit: github.Bool(true)})
	if err != nil {
		fmt.Println("err")
	}

	fmt.Println(repo)
	fmt.Println(res)
}