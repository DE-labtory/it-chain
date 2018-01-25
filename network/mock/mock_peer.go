package mock

import (
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "it-chain/network/protos"
	"io"
	"fmt"
	"golang.org/x/net/context"
)