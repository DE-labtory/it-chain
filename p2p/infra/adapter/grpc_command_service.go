package adapter

import (
	"time"

	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyPeerId = errors.New("empty nodeid proposed")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection ")

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

// message dispatcher sends messages to other nodes in p2p network
type GrpcCommandService struct {
	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *GrpcCommandService {
	return &GrpcCommandService{
		publish: publish,
	}
}

//request leader information in p2p network to the node specified by peerId
func (gcs *GrpcCommandService) RequestLeaderInfo(connectionId string) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	//message deliver command for delivering leader info
	deliverCommand, err := CreateGrpcDeliverCommand("LeaderInfoRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, connectionId)

	return gcs.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (gcs *GrpcCommandService) RequestPeerList(peerId p2p.PeerId) error {

	if peerId.Id == "" {
		return ErrEmptyPeerId
	}
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateGrpcDeliverCommand("PeerListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return gcs.publish("Command", "message.deliver", deliverCommand)
}

func (gcs *GrpcCommandService) DeliverLeaderInfo(connectionId string, leader p2p.Leader) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	deliverCommand, err := CreateGrpcDeliverCommand("UpdateLeader", leader)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, connectionId)

	return gcs.publish("Command", "message.deliver", deliverCommand)
}

func (gcs *GrpcCommandService) DeliverPeerLeaderTable(connectionId string, peerLeaderTable p2p.PeerLeaderTable) error {

	if connectionId== "" {
		return ErrEmptyPeerId
	}

	if len(peerLeaderTable.PeerList) ==0 {
		return p2p.ErrEmptyPeerList
	}

	//create peer table message
	peerLeaderTableMessage := p2p.PeerLeaderTableMessage{
		PeerLeaderTable:peerLeaderTable,
	}

	grpcDeliverCommand, err := CreateGrpcDeliverCommand("PeerTableDeliver", peerLeaderTableMessage)

	if err != nil {
		return err
	}

	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	return gcs.publish("Command", "message.deliver", grpcDeliverCommand)
}

func (gcs *GrpcCommandService) DeliverRequestVoteMessages(connectionIds ...string) error{

	requestVoteMessage := p2p.RequestVoteMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("PeerTableDeliver", requestVoteMessage)

	for _, connectionId := range connectionIds{
		grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)
	}

	gcs.publish("Command", "message.send", grpcDeliverCommand)

	return nil

}

func (gcs *GrpcCommandService) DeliverVoteLeaderMessage(connectionId string) error {
	voteLeaderMessage := p2p.VoteLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("VoteLeaderProtocol", voteLeaderMessage)

	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	gcs.publish("Command", "message.send", grpcDeliverCommand)

	return nil

}

func (gcs *GrpcCommandService) DeliverUpdateLeaderMessage(connectionId string, peer p2p.Peer) error {

	updateLeaderMessage := p2p.UpdateLeaderMessage{}

	grpcDeliverCommand, _ := CreateGrpcDeliverCommand("UpdateLeaderProtocol", updateLeaderMessage)

	gcs.publish("Command", "message.deliver", grpcDeliverCommand)

	return nil

}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (p2p.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return p2p.GrpcDeliverCommand{}, err
	}

	return p2p.GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
