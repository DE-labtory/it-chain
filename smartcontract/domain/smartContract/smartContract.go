package smartContract

type Version struct {
}

type SmartContract struct {
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
}

func NewSmartContract(repositoryName string, gitUrl string, path string, commitHash string) *SmartContract {

	return &SmartContract{
		RepositoryName: repositoryName,
		CommitHash:     commitHash,
		GitUrl:         gitUrl,
		Path:           path,
	}
}
