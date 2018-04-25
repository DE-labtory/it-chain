package model

type ICodeConfiguration struct {
	RepositoryPath string
}

func NewIcodeConfiguration() ICodeConfiguration {
	return ICodeConfiguration{
		RepositoryPath: "empty",
	}
}
