package service

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

var username string
var password string

var homepage = "https://github.com"

var client *github.Client

//todo read from config
//todo private public from config
func init() {

	username = "steve@buzzni.com"
	password = "itchain123"

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client = github.NewClient(tp.Client())
}

func GetCurrentGithubUserName() string {
	return username
}

func CreateRepository(name string) (*github.Repository, *github.Response, error) {

	ctx := context.Background()
	repo, res, err := client.Repositories.Create(ctx, "", &github.Repository{
		Name:     &name,
		Private:  &[]bool{false}[0],
		Homepage: &homepage,
	})

	if err != nil {
		return repo, res, err
	}

	return repo, res, nil
}
