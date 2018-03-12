package main

import (
	"os/exec"
	"errors"
)

func CommitAll(dir string, comment string) (error) {
	cmd_add := exec.Command("git", "add", ".")
	cmd_add.Dir = dir
	error := cmd_add.Run()
	if error != nil {
		return errors.New("Add Error")
	}

	//out, err := exec.Command("git", "commit", "-m", comment).Output()
	//fmt.Println(string(out))
	//if err != nil {
	//	return errors.New(err.Error())
	//}

	cmd_commit := exec.Command("git", "commit", "-m", comment)
	cmd_commit.Dir = dir
	error = cmd_commit.Run()
	if error != nil {
		return errors.New("Commit Error")
	}

	return nil
}

func PushRepos(dir string) (error) {
	cmd := exec.Command("git", "push")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return errors.New("Push Error")
	}

	return nil
}

func CommitAndPush(dir string, comment string) (error) {
	error := CommitAll(dir, comment)
	if error != nil {
		return error
	}
	error = PushRepos(dir)
	if error != nil {
		return error
	}
	return nil
}

func main() {
	err := CommitAndPush("/Users/hackurity/Documents/it-chain/test/bloom", "It-Chain Smart Contract \"junbeomlee_bloom\" Deploy")
	if err != nil {
		panic(err.Error())
	}
}