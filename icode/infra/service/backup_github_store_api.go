package service

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var homepage = "https://github.com"

var client *github.Client

//todo read from config
//todo private public from config
type BackupGithubStoreApi struct {
	client      *github.Client
	homepageUrl string
	storename   string
	username    string
	password    string
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
		storename:   user.GetLogin(),
		username:    username,
		password:    password,
	}, nil
}

func (bgs BackupGithubStoreApi) GetName() string {
	return bgs.username
}

func (bgs BackupGithubStoreApi) GetHomepageUrl() string {
	return bgs.homepageUrl
}

//create backup repo
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

//push to backup server
func (bgs BackupGithubStoreApi) PushRepository(repositoryPath string) error {

	au := &http.BasicAuth{Username: bgs.username, Password: bgs.password}

	r, err := git.PlainOpen(repositoryPath)

	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       au,
	})

	if err != nil {
		return err
	}

	return nil
}
