package service

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

var homepage = "https://github.com"

var client *github.Client

//todo read from config
//todo private public from config

type BackupGithubStoreApi struct {
	client      *github.Client
	homepageUrl string
	name        string
}

func NewBackupGithubStoreApi(username string, password string) (*BackupGithubStoreApi, error) {

	if username == "" {
		username = "steve@buzzni.com"
	}

	if password == "" {
		password = "itchain123"
	}

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client = github.NewClient(tp.Client())

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")

	if err != nil {
		return nil, err
	}

	return &BackupGithubStoreApi{
		client:      client,
		homepageUrl: user.GetHTMLURL(),
		name:        user.GetLogin(),
	}, nil
}

func (bgs BackupGithubStoreApi) GetName() string {
	return bgs.name
}

func (bgs BackupGithubStoreApi) GetHomepageUrl() string {
	return bgs.homepageUrl
}

func (bgs BackupGithubStoreApi) CreateRepository(name string) (*github.Repository, error) {

	ctx := context.Background()
	r, _, err := client.Repositories.Create(ctx, "", &github.Repository{
		Name:     &name,
		Private:  &[]bool{false}[0],
		Homepage: &homepage,
	})

	if err != nil {
		return r, err
	}

	return r, nil
}

func (bgs BackupGithubStoreApi) GetRepositoryList() []string {
	//todo
	return nil
}
