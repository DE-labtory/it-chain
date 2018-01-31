package main

import (
	"io"
	"os"
	"context"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
)

func PullImage(imageName string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
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