package webhook

import (
	"errors"
	pb "it-chain/service/webhook/proto"
	"context"
)

type WebhookServiceImpl struct {

	urls []string

}

func NewWebhookService() (WebhookService, error) {

	urls := make([]string, 0)

	// 만약 DB 사용시 로드 하여 urls에 로드하기
	wi := &WebhookServiceImpl {
		urls: urls,
	}

	return wi, nil

}

func (wi *WebhookServiceImpl) Register(ctx context.Context, in *pb.WebhookRequest) (*pb.WebhookResponse, error) {

	if len(in.Urls) == 0 {
		return &pb.WebhookResponse{"INVALID URL"}, nil
	}

	wi.urls = append(wi.urls, in.Urls)

	return &pb.WebhookResponse{"SUCCESS TO REGISTER WEBHOOK URL"}, nil

}