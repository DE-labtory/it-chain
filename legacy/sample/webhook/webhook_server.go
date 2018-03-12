package main

import (
	"github.com/it-chain/it-chain-Engine/legacy/service/webhook"
	"log"
	"github.com/it-chain/it-chain-Engine/legacy/domain"
	"time"
)

func main() {

	webhookService, err := webhook.NewWebhookService()
	if err != nil {
		log.Fatalf("failed to initialize webhook service : %v", err)
	}

	confirmedBlock := &domain.Block{}

	go func() {
		for {
			log.Println("[BLOCK CONFIRMED]")
			err = webhookService.SendConfirmedBlock(confirmedBlock)
			if err != nil {
				log.Fatalf("failed to send confirmed block : %v", err)
			}

			time.Sleep(3000 * time.Millisecond)
		}
	}()

	if err = webhookService.Serve(50070); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}

}