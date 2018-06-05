FROM golang:1.9
ADD . /go/src/github.com/it-chain/it-chain-Engine
WORKDIR /go/src/github.com/it-chain/it-chain-Engine
ENV GOBIN=/go/bin
CMD ["./start.sh"]