package comm

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

const defaultTimeout = time.Second * 3

//get time from config
//timeOut := viper.GetInt("grpc.timeout")
func NewConnectionWithAddress(peerAddress string,  tslEnabled bool, creds credentials.TransportCredentials) (*grpc.ClientConn, error){

	var opts []grpc.DialOption

	if tslEnabled {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithTimeout(defaultTimeout))

	conn, err := grpc.Dial(peerAddress, opts...)
	if err != nil {
		return nil, err
	}

	return conn, err
}