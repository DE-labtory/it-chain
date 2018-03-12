package webhook

import "github.com/it-chain/it-chain-Engine/domain"

type WebhookService interface {

	SendConfirmedBlock(block *domain.Block) error

	Serve(port int) error

}
