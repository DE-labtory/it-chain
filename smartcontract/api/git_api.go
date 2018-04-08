package api

import (
	"fmt"
	"os/user"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var tmp = "./.tmp"

type GitApi struct {
}

func Clone(url string) error {

	currentUser, err := user.Current()
	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")

	if err != nil {
		return err
	}

	r, err := git.PlainClone(tmp+"/"+url, false, &git.CloneOptions{
		URL:               url,
		Auth:              sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	committer, err := r.CommitObjects()

	lastCommit, err := committer.Next()
	commitHash := lastCommit.Hash.String()

	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Print(commitHash)

	return nil
}
