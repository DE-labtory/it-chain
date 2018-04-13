package icodeMeta

type Version struct {
}

type ICodeMeta struct {
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func NewItCode(repositoryName string, gitUrl string, path string, commitHash string) *ICodeMeta {

	return &ICodeMeta{
		RepositoryName: repositoryName,
		CommitHash:     commitHash,
		GitUrl:         gitUrl,
		Path:           path,
	}
}
