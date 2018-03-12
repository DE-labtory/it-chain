package main

import (
	"os"
	"context"
	"fmt"
	"io"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	//	"docker.io/go-docker/api/types/volume"
	"strings"
)

func main() {
	/*** Docker Set ***/
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		fmt.Println("error NewEnvClient")
		return
	}

	imageName := "docker.io/library/golang:1.9.2-alpine3.6"
	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		fmt.Println("error ImagePull")
		return
	}
	io.Copy(os.Stdout, out)

	/*** Volume Test ***/
	GOPATH := os.Getenv("GOPATH")
	fmt.Println(GOPATH)
	createContHostConfig := container.HostConfig{
		Binds:           []string{GOPATH + "/src/it-chain:/go/src/it-chain",},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
	}, &createContHostConfig, nil, "")
	if err != nil {
		fmt.Println("error ContainerCreate")
		return
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("error ContainerStart")
		return
	}

}
