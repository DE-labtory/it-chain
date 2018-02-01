package main

import (
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"errors"
	"net/http"
)

type GithubResponse struct {
	Message	  string
	Id        int			`json:"id"`
	Name      string		`json:"name"`
	Full_name string		`json:"full_name"`
	Html_url  string		`json:"html_url"`
	Owner     struct {
		Login	string		`json:"login"`
	}						`json:"owner"`
	Source    struct {
		Id        int		`json:"id"`
		Name      string	`json:"name"`
		Full_name string	`json:"full_name"`
		Html_url  string	`json:"html_url"`
		Owner     struct {
			Login	string	`json:"login"`
		}					`json:"owner"`
	}						`json:"source"`
}

type GithubRequestCreateRepos struct {
	Name 		string		`json:"name"`
	Description	string		`json:"description"`
}

func CreateRepos(repos_name string, token string) (GithubResponse, error) {
	body := GithubResponse{}
	api_url := "https://api.github.com/user/repos?access_token=" + token

	param := GithubRequestCreateRepos{repos_name, repos_name}
	param_string, err := json.Marshal(param)
	if err != nil {
		return body, errors.New("Marshal Error")
	}

	fmt.Println(param)
	fmt.Println(param_string)
	fmt.Println(bytes.NewBuffer(param_string))
	res, err := http.Post(api_url, "application/json", bytes.NewBuffer([]byte(param_string)))
	if err != nil {
		return body, errors.New("Error Create Repos")
	}
	defer res.Body.Close()

	if res.Header.Get("Status") != "201 Created" {
		return body, errors.New(res.Header.Get("Status"))
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return body, errors.New("Unmarshal Error")
	}

	return body, nil
}

func main() {
	body, err := CreateRepos("jun_lee", "1f8f0f1e16bb4b98e3d5d6113b34a3d8cb6ab8d2")
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(body)
}