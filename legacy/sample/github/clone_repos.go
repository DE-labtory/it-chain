package main

import (
	"os/exec"
)

func CloneRepos(repos_path string, dir string) (error) {
	cmd := exec.Command("git", "clone", "https://github.com/"+repos_path+".git")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
	}

	return nil
}

func main() {
	err := CloneRepos("hackurity01/bloom2", "/Users/hackurity/Documents/it-chain/test")
	println(err)
}