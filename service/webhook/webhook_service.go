package webhook

import "it-chain/domain"

type WebhookService interface {

	SendConfirmedBlock(block *domain.Block) (error)

}
