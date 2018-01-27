package smartcontract

import (
	"errors"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"os/exec"
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

type GithubResponseCommits struct {
	Message	  string
	Sha       string		`json:"sha"`
	Committer struct {
		Login	string		`json:"login"`
	}						`json:"committer"`
}

type GithubRequestCreateRepos struct {
	Name 		string		`json:"name"`
	Description	string		`json:"description"`
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

func CreateRepos(repos_name string, token string) (GithubResponse, error) {
	var body = GithubResponse{}
	api_url := "https://api.github.com/user/repos?access_token=" + token

	param := GithubRequestCreateRepos{repos_name, repos_name}
	param_bytes, err := json.Marshal(param)
	if err != nil {
		return body, errors.New("Marshal Error")
	}

	res, err := http.Post(api_url, "application/json", bytes.NewBuffer(param_bytes))
	if err != nil {
		return body, errors.New("Error Create Repos")
	}
	defer res.Body.Close()

	if res.Header.Get("Status") != "201 Created" {
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

func CloneRepos(repos_path string, dir string) (error) {
	cmd := exec.Command("git", "clone", "https://github.com/"+repos_path+".git")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
	}

	return nil
}

func ChangeRemote(repos_path string, dir string) (error) {
	cmd := exec.Command("git", "remote", "set-url", "origin", "https://github.com/"+repos_path+".git")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
	}

	return nil
}

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

func PushRepos(dir string) (error) {
	cmd := exec.Command("git", "push")
	cmd.Dir = dir
	error := cmd.Run()
	if error != nil {
		return error
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