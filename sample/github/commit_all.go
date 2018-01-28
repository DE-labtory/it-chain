package main

import (
	"os/exec"
)


func CommitAll(dir string, comment string) (error) {
	cmd_add := exec.Command("git", "add", ".")
	cmd_add.Dir = dir
	error := cmd_add.Run()
	if error != nil {
		return error
	}

	cmd_commit := exec.Command("git", "commit", "-m", comment)
	cmd_commit.Dir = dir
	error = cmd_commit.Run()
	if error != nil {
		return error
	}

	return nil
}

func main() {
	err := CommitAll("/Users/hackurity/Documents/it-chain/test/bloom2", "test")
	if err != nil {
		panic(err)
	}
}
