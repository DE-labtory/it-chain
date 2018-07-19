FROM golang:1.9
ADD . /go/src/github.com/it-chain/engine
WORKDIR /go/src/github.com/it-chain/engine
ENV GOBIN=/go/bin
CMD ["./start.sh"]