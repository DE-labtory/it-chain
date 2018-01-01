#  Go GRPC 사용법



## Prerequisite

1. Go 1.6 이상

2. Protoc

   - mac: brew install protobuf

3. Protoc-gen-go

   - go get github.com/golang/protobuf/protoc-gen-go

   - $GOPATH/bin에 있어야함!!

     ​

   ​

   ​

## Proto 파일 생성

```
syntax = "proto3";
package proto;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  rpc SayHelloAgain (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

##### service는 request와 return을 정해주며 message는 request와 return의 형식을 정해준다.



## sample.pb.go 생성 방법

```
protoc -I ./ sample.proto --go_out=plugins=grpc:./
# protoc -I [src위치] [sample.proto 위치] --go_out=plugins=grpc:[output 위치]
```



##### sample.pb.go가 같은 path에 생성



## 클라이언트, 서버 코드

### 클라이언트

```
package main

import (
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "it-chain/sample/grpc/proto"
	"os"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	r, err = c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
```



### 서버

```
package main


import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "it-chain/sample/grpc/proto"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```





## 실행

서버 실행

```
go run grpc.go
```



클라이언트 실행

```
go run grpc_client.go
2018/01/01 16:30:44 Greeting: Hello world
2018/01/01 16:30:44 Greeting: Hello again world
```


## 참고 문서
https://grpc.io/docs/quickstart/go.html

https://developers.google.com/protocol-buffers/docs/gotutorial


