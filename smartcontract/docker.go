package smartcontract

import (
	"io"
	"os"
	"io/ioutil"
	"fmt"
	"bytes"
	"context"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"strings"
)

// imageName := "docker.io/library/node"
/*
func PullDockerImage(imageName string) {
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




}

func CreateDockerContainer(image string) {

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{"/bin/bash"},
		Tty: true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

}

func StartDockerContainer() {

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func CopyTarToDockerContainer() {

	file, err := ioutil.ReadFile("/Users/hackurity/go/src/docker_client_test/kk.tar")
	if err != nil {
		fmt.Print(err)
	}

	err = cli.CopyToContainer(ctx, resp.ID, "/home/node", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		panic(err)
	}
}
*/


// Temporary Function. It have to be splited for legibility and maintenance.
func PullAndCopyAndRunDocker(imageName string, tarPath string) *client.Client{
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

	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{"/bin/bash"},
		Tty: true,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadFile(tarPath)
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println(file)



	err = cli.CopyToContainer(ctx, resp.ID, "/home", bytes.NewReader(file), types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	})
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	return cli
}