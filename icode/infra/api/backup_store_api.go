package api

import "github.com/google/go-github/github"

type BackupStoreApi interface {
	CreateRepository(name string) (*github.Repository, error)
	GetRepositoryList() []string
	GetHomepageUrl() string
	GetName() string
	PushRepository(repositoryPath string) error
}
