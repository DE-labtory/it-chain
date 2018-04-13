package itcode

type Version struct {
}

type ItCode struct {
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func NewItCode(repositoryName string, gitUrl string, path string, commitHash string) *ItCode {

	return &ItCode{
		RepositoryName: repositoryName,
		CommitHash:     commitHash,
		GitUrl:         gitUrl,
		Path:           path,
	}
}
