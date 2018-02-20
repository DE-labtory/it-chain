package main

import (
	"log"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"context"
	"fmt"
)

func main(){

	address := "127.0.0.1:5555"

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	transaction, err := c.PostTransaction(context.Background(), &pb.TxData{
		Jsonrpc: "json2.0",
		Method: pb.TxData_Invoke,
		ContractID: "zxcjzixcj",
		Params: &pb.Params{},
	})

	fmt.Print(transaction)

	if err != nil {
		log.Println("could not greet: %v", err)
	}
}
