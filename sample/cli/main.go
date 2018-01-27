package main

import (
	"github.com/urfave/cli"
	"os"
	"fmt"
)

func main(){

	//실행 방법
	//1. go install로 bin으로 만들어서 사용
	//2. go run으로 바로 사용

	//사용법
	//go run main.go --ip [IP ADDRESS], -i [IP ADDRESS] --peerid [PEERID], --p [PEERID]
	//로 실행시 ipAddress, peerID에 각각 값이 들어간다.

	var ipAddress = ""
	var peerID = ""

	app := cli.NewApp()
	app.Name = "PEER"
	app.Usage = "fight the loneliness!"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:        "ip, i",
			Usage:       "IP Address of Boot Peer",
			Destination: &ipAddress,
		},
		cli.StringFlag{
			Name:        "peerid, p",
			Usage:       "ID of Peer",
			Destination: &peerID,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("exec")

		if ipAddress == ""{
			fmt.Println("no ip address")
			return nil
		}

		fmt.Println(ipAddress)
		fmt.Println(peerID)
		return nil
	}

	app.Run(os.Args)
}