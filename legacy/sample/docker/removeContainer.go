package main

import (
	//"fmt"
	"os/exec"
	"bytes"
	"io"
)

func main() {
	 //docker rm $(docker ps -a -f "ancestor=golang:1.9.2-alpine3.6" -q)
	 //docker ps -a -f "ancestor=golang:1.9.2-alpine3.6" -q | xargs -I {} docker rm {}
	//cmd := exec.Command("docker", "ps", "-a", "-f", "ancestor=golang:1.9.2-alpine3.6", "-q")
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	fmt.Println("test docker container remove error")
	//}

	c1 := exec.Command("docker", "ps", "-a", "-f", "ancestor=golang:1.9.2-alpine3.6", "-q")
	c2 := exec.Command("xargs", "-I", "{}", "docker", "rm", "{}")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
}
