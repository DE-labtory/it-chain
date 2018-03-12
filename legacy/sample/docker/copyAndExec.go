package main

import (
	"fmt"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/it-chain/it-chain-Engine/legacy/smartcontract"
	"bufio"
	"os"
	"strings"
	"bytes"
	"context"
	"io"
	"io/ioutil"
)

func main() {
	err := smartcontract.MakeTar("/tmp/test", "/tmp")
	if err != nil {
		fmt.Println(err.Error())
	}
	imageName := "docker.io/library/golang:1.9.2-alpine3.6"
	tarPath := "/tmp/test.tar"

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
	fmt.Println("Passed ImagePull")

	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		//Cmd: []string{"/bin/sh"},
		Cmd: []string{"/go/src/test"},
		Tty: true,
		AttachStdout: true,
		AttachStderr: true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed ContainerCreate")

	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		fmt.Print(err)
	}


	err = cli.CopyToContainer(ctx, resp.ID, "/go/src/", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed CopyToContainer Go File")


	fmt.Println("============================")
	fmt.Println("resp.ID : " + resp.ID)
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Passed ContainerStart")


	/* get docker output
	----------------------*/
	fmt.Println("=============<Docker Output>===============")
	reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}