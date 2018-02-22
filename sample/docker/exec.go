package main

import (
	"fmt"
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
)

func main() {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	r, err := cli.ContainerExecCreate(ctx, "031948a18830", types.ExecConfig{
		Cmd: []string{"go", "build", "-o", "test1", "/go/src/fileio.go"},
		User: "root",
		//Cmd: []string{"touch","/home/aa"},
	})
	if err != nil {
		panic(err)
	}

	err = cli.ContainerExecStart(ctx, r.ID, types.ExecStartCheck{
		Detach: true,
		Tty:    true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(r)
}