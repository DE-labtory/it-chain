package main

import (
	"os/exec"
)

func PushRepos(dir string) (error) {
	cmd := exec.Command("git", "push")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
	}

	return nil
}

func main() {
	err := PushRepos("/Users/hackurity/Documents/it-chain/test/bloom")
	println(err)
}