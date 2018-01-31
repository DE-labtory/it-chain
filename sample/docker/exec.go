package main

import (
	"fmt"
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	r, err := cli.ContainerExecCreate(ctx, "031948a18830", types.ExecConfig{
		Cmd: []string{"go", "run", "/go/src/fileio.go"},
		User: "root",
		WorkingDir: "/go/src",
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