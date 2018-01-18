package chaincode

import (
	"errors"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

type GithubResponse struct {
	Id        int			`json: id`
	Name      string		`json: name`
	Full_name string		`json: full_name`
	Html_url  string		`json: html_url`
	Owner     struct {
		Login	string		`json: login`
	}						`json: owner`
	Source    struct {
		Id        int		`json: id`
		Name      string	`json: name`
		Full_name string	`json: full_name`
		Html_url  string	`json: html_url`
		Owner     struct {
			Login	string	`json: login`
		}					`json: owner`
	}						`json: source`
}

func GetRepos(repos_path string) (GithubResponse, error) {
	var body = GithubResponse{}
	api_url := "https://api.github.com/repos/" + repos_path

	res, err := http.Get(api_url)
	if err != nil {
		return body, errors.New("Wrong repository path")
	}
	defer res.Body.Close()

	if res.Header.Get("Status") != "200 OK" {
		return body, errors.New("Not Found")
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return body, errors.New("Unmarshal Error")
	}

	return body, nil
}

func ForkRepos(repos_path string, token string) (GithubResponse, error) {
	var body = GithubResponse{}
	api_url := "https://api.github.com/repos/" + repos_path + "/forks?access_token=" + token

	res, err := http.Post(api_url, "", bytes.NewBufferString(""))
	if err != nil {
		return body, errors.New("Wrong repository path")
	}
	defer res.Body.Close()

	if res.Header.Get("Status") != "202 Accepted" {
		return body, errors.New("Post Error")
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return body, errors.New("Unmarshal Error")
	}

	return body, nil
}
