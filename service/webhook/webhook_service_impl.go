package webhook

import (
	"errors"
	pb "it-chain/service/webhook/proto"
	"context"
	"it-chain/domain"
	"encoding/json"
	"net/url"
	"bytes"
	"net/http"
	"strings"
	"google.golang.org/grpc"
	"net"
	"log"
	"google.golang.org/grpc/reflection"
	"fmt"
)

type WebhookServiceImpl struct {

	infos []webhookInfo

}

type webhookInfo struct {

	payloadURL url.URL

}

func NewWebhookService() (WebhookService, error) {

	infos := make([]webhookInfo, 0)

	wi := &WebhookServiceImpl {
		infos: infos,
	}

	return wi, nil

}

func (wi *WebhookServiceImpl) Register(ctx context.Context, in *pb.WebhookRequest) (*pb.WebhookResponse, error) {

	parseURL, valid := urlValidCheck(in.PayloadURL)
	if !valid {
		return &pb.WebhookResponse{"Invalid URL"}, errors.New("Invalid URL")
	}

	_, isExist := wi.urlExistCheck(parseURL)
	if isExist {
		return &pb.WebhookResponse{"Already exists"}, errors.New("Already exist")
	}

	log.Printf("REGISTED URL : %s", parseURL.String())

	wi.infos = append(wi.infos, webhookInfo{*parseURL})
	return &pb.WebhookResponse{"Success to register your webhook"}, nil

}

func (wi *WebhookServiceImpl) Remove(ctx context.Context, in *pb.WebhookRequest) (*pb.WebhookResponse, error) {

	parseURL, valid := urlValidCheck(in.PayloadURL)
	if !valid {
		return &pb.WebhookResponse{"Invalid URL"}, errors.New("Invalid URL")
	}

	index, isExist := wi.urlExistCheck(parseURL)
	if isExist {
		// Remove duplicated element
		wi.infos = append(wi.infos[:index], wi.infos[index+1:]...)
		return &pb.WebhookResponse{"Success to remove your webhook"}, nil
	}

	return &pb.WebhookResponse{"Requested URL is not exist"}, errors.New("Requested URL is not exist")

}

func (wi *WebhookServiceImpl) Serve(port int) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	defer lis.Close()

	grpcServer := grpc.NewServer()

	pb.RegisterWebhookServer(grpcServer, wi)

	reflection.Register(grpcServer)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to launch grpc server : %v", err)
		grpcServer.Stop()
		return err
	}

	defer grpcServer.Stop()

	return nil

}

func (wi *WebhookServiceImpl) SendConfirmedBlock(block *domain.Block) (error) {

	if block == nil {
		return errors.New("block should not be nil")
	}

	if len(wi.infos) == 0 {
		log.Printf("There is no payload URL")
		return nil
	}

	blockBytes, err := json.Marshal(block)
	if err != nil {
		return errors.New("An error is occured during converting process")
	}

	buff := bytes.NewBuffer(blockBytes)

	go func() {
		for _, info := range wi.infos {
			log.Printf("PAYLOAD URL : %s", info.payloadURL.String())
			_, err := http.Post(info.payloadURL.String(), "application/json", buff)
			if err != nil {
				log.Fatalf("An error during the sending process : %v", err)
			}
		}
	}()

	return nil

}

func (wi *WebhookServiceImpl) urlExistCheck(url *url.URL) (int, bool) {

	for idx, info := range wi.infos {
		if strings.Compare(info.payloadURL.String(), url.String()) == 0 {
			return idx, true
		}
	}

	return -1, false

}

func urlValidCheck(rawURL string) (*url.URL, bool) {

	if len(rawURL) == 0 {
		return nil, false
	}

	parseURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, false
	}

	return parseURL, true

}