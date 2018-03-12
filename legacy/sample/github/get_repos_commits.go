package main

import (
	"io/ioutil"
	"encoding/json"
	"errors"
	"net/http"
	"fmt"
)

type GithubResponseCommits struct {
	Message	  string
	Sha       string			`json: sha`
	Committer struct {
		Login	string		`json: login`
	}						`json: committer`
}

func GetReposCommits(repos_path string) ([]GithubResponseCommits, error) {
	var body []GithubResponseCommits
	api_url := "https://api.github.com/repos/" + repos_path + "/commits"

	res, err := http.Get(api_url)
	if err != nil {
		return body, errors.New("Wrong repository path")
	}
	defer res.Body.Close()

	if res.Header.Get("Status") != "200 OK" {
		return body, errors.New("Not Found")
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return body, errors.New("ioutil Error")
	}

	err = json.Unmarshal(bodyBytes, &body)
	return body, nil
}

func main() {
	git, err := GetReposCommits("hackurity01/test")
	if err != nil {
		panic(err)
	}
	fmt.Print(git)
}