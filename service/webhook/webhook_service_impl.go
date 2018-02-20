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
)

type WebhookServiceImpl struct {

	urls []url.URL

}

func NewWebhookService() (WebhookService, error) {

	urls := make([]url.URL, 0)

	// 만약 DB 사용시 로드 하여 urls에 로드하기
	wi := &WebhookServiceImpl {
		urls: urls,
	}

	return wi, nil

}

func (wi *WebhookServiceImpl) Register(ctx context.Context, in *pb.WebhookRequest) (*pb.WebhookResponse, error) {

	url, valid := urlValidCheck(in.Urls)
	if !valid {
		return &pb.WebhookResponse{"INVALID URL"}, errors.New("INVALID URL")
	}

	// URL 중복 체크 필요
	wi.urls = append(wi.urls, *url)

	return &pb.WebhookResponse{"SUCCESS TO REGISTER YOUR WEBHOOK URL"}, nil

}

func (wi *WebhookServiceImpl) Remove(ctx context.Context, in *pb.WebhookRequest) (*pb.WebhookResponse, error) {

	url, valid := urlValidCheck(in.Urls)
	if !valid {
		return &pb.WebhookResponse{"INVALID URL"}, errors.New("INVALID URL")
	}

	// 해당 slice에서 매칭되는 URL 찾아서 지우기

	return &pb.WebhookResponse{"SUCCESS TO REMOVE YOUR WEBHOOK URL"}, nil

}

func (wi *WebhookServiceImpl) SendConfirmedBlock(block *domain.Block) (error) {

	if block == nil {
		return errors.New("block should not be nil")
	}

	blockBytes, err := json.Marshal(block)
	if err != nil {
		return errors.New("An error is occured during converting process")
	}

	buff := bytes.NewBuffer(blockBytes)

	// json, xml? 등의 방식으로 보낼 때 뭐로 보낼지 정할 수 있도록 구현하는게 좋을 것 같다. ( 전송 타입 ? )
	go func() {
		for _, url := range wi.urls {
			res, err := http.Post(url.String(), "application/json", buff)
		}
	}()

	return nil

}

func urlValidCheck(rawURL string) (*url.URL, bool) {

	if len(rawURL) == 0 {
		return nil, false
	}

	parseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, false
	}

	return parseURL, true

}