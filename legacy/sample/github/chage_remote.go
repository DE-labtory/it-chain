package main

import (
	"os/exec"
)

func ChangeRemote(repos_path string, dir string) (error) {
	cmd := exec.Command("git", "remote", "set-url", "origin", "https://github.com/"+repos_path+".git")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
	}

	return nil
}

func main() {
	err := ChangeRemote("hackurity01/createTest", "/Users/hackurity/Documents/it-chain/test/bloom2")
	println(err)
}