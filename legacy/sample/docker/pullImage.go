package main

import (
	"io"
	"os"
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

func PullImage(imageName string) error {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	return nil
}

func main() {
	PullImage("docker.io/library/golang:rc-alpine")
}