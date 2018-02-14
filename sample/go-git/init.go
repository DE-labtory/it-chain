package main

import (
"os"
"gopkg.in/src-d/go-git.v4"
"fmt"
)

func main() {
	//Info("git clone https://github.com/src-d/go-git")

	_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      "https://github.com/src-d/go-git",
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("test")
	}

	//CheckIfError(err)
}