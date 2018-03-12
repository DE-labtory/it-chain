package service

import (
	"time"
	"github.com/it-chain/it-chain-Engine/legacy/common"
	"math"
	"github.com/it-chain/it-chain-Engine/legacy/domain"
	"github.com/it-chain/it-chain-Engine/legacy/network/comm"
	"sync"
	"github.com/rs/xid"
	pb "github.com/it-chain/it-chain-Engine/legacy/network/protos"
	"github.com/it-chain/it-chain-Engine/legacy/network/comm/msg"
	"github.com/spf13/viper"
	"strconv"
	"github.com/it-chain/it-chain-Engine/legacy/service/webhook"
)

var logger_pbftservice = common.GetLogger("pbft_service")

//todo peerID를 어디서 가져올 것인가??
type PBFTConsensusService struct {
	consensusStates      map[string]*domain.ConsensusState
	comm                 comm.ConnectionManager
	identity             *domain.Peer
	peerService          PeerService
	blockService         BlockService
	smartContractService SmartContractService
	webHookService       webhook.WebhookService
	transactionService   TransactionService
	sync.RWMutex
}

func NewPBFTConsensusService(comm comm.ConnectionManager,webHookService webhook.WebhookService, peerService PeerService, blockService BlockService,identity *domain.Peer, smartContractService SmartContractService, transactionService TransactionService) ConsensusService{

	pbft := &PBFTConsensusService{
		consensusStates: make(map[string]*domain.ConsensusState),
		comm:comm,
		peerService: peerService,
		blockService: blockService,
		smartContractService: smartContractService,
		transactionService: transactionService,
		webHookService: webHookService,
		identity: identity,
	}

	i, _ := strconv.Atoi(viper.GetString("consensus.batchTime"))

	broadCastPeerTableBatcher := NewBatchService(time.Duration(i)*time.Second,pbft.startConsensus,false)
	broadCastPeerTableBatcher.Add("start consensus")
	broadCastPeerTableBatcher.Start()

	comm.Subscribe("Handle consensus msg",pbft.ReceiveConsensusMessage)

	return pbft
}

//tested
//Consensus 시작
//만약 합의에 들어가는 peerID가 없다면 바로 block에 저장
//1. Consensus의 state를 추가한다.
//2. 합의할 block을 consensusMessage에 담고 prepreMsg로 전파한다.
func (cs *PBFTConsensusService) StartConsensus(view *domain.View, block *domain.Block){

	cs.Lock()
	consensusState := domain.NewConsensusState(view,xid.New().String(),block,domain.PrePrepared,cs.EndConsensusState,300)
	cs.consensusStates[consensusState.ID] = consensusState

	sequenceID := time.Now().UnixNano()
	preprepareConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,*view,sequenceID,consensusState.Block,cs.identity.PeerID,domain.PreprepareMsg)

	cs.broadcastMessage(preprepareConsensusMessage)

	consensusState.CurrentStage = domain.Prepared
	cs.Unlock()
}

func (cs *PBFTConsensusService) startConsensus(interface{}){

	//1. 혼자인경우
	//2. leader아닌경우

	if (cs.peerService.GetLeader() == nil || cs.identity.PeerID != cs.peerService.GetLeader().PeerID) && len(cs.peerService.GetPeerTable().GetPeerList()) != 0{
		common.Log.Println("Not leader")
		return
	}

	common.Log.Println("start Consesnsus")

	transactions, err := cs.transactionService.GetTransactions(100)

	if err !=nil{
		common.Log.Error("Fail to get transactions",err.Error())
		return
	}

	if len(transactions) == 0{
		common.Log.Println("No tx to consesnsus")
		return
	}

	for i := 0; i < len(transactions); i++ {
		cs.smartContractService.ValidateTransaction(transactions[i])
		transactions[i].GenerateHash()
	}

	cs.transactionService.DeleteTransactions(transactions)
	block, err := cs.blockService.CreateBlock(transactions,cs.identity.PeerID)

	if err != nil{
		common.Log.Error("Fail to create block",err.Error())
		return
	}

	//todo 혼자인경우
	//common.Log.Error(cs.peerService.GetPeerTable())
	if len(cs.peerService.GetPeerTable().GetPeerList()) == 0{

		flag, err := cs.blockService.VerifyBlock(block)

		if err != nil{
			common.Log.Error("Verify block error:",err.Error())
		}

		if flag{
			common.Log.Println("Add block")
			asd, err := cs.blockService.AddBlock(block)

			common.Log.Println(asd)

			if err !=nil{
				common.Log.Println(err.Error())
			}
			cs.webHookService.SendConfirmedBlock(block)
			common.Log.Println("Add block2")
		}
		return
	}

	peerIDs := make([]string,0)

	for _, peer := range cs.peerService.GetPeerTable().GetAllPeerList(){
		peerIDs = append(peerIDs, peer.PeerID)
	}

	view := &domain.View{
		ID: xid.New().String(),
		LeaderID: cs.identity.PeerID,
		PeerID: peerIDs,
	}

	cs.StartConsensus(view,block)
}

func (cs *PBFTConsensusService) GetCurrentConsensusState() map[string]*domain.ConsensusState{
	return cs.consensusStates
}

func (cs *PBFTConsensusService) StopConsensus(){

	cs.Lock()

	defer cs.Unlock()

	for consensusState := range cs.consensusStates {
		cs.consensusStates[consensusState].End()
		delete(cs.consensusStates, consensusState)
	}
}

//consensusMessage가 들어옴
//todo FromConsensusProtoMessage에서 block변환도 해야함
//todo time을 config로 부터 읽어야함
//todo 다음 block이 먼저 들어올 경우 고려해야함,
//todo 블록의 높이와 이전 블록 해시가 올바른지 확인
func (cs *PBFTConsensusService) ReceiveConsensusMessage(msg msg.OutterMessage){

	if consensusMsg := msg.Message.GetConsensusMessage(); consensusMsg ==nil{
		return
	}

	common.Log.Println("Received Consensus Msg")

	message := msg.Message
	cm:= message.GetConsensusMessage()

	if cm == nil{
		logger_pbftservice.Errorln("consensus Message is empty")
		return
	}

	consensusMessage := domain.FromConsensusProtoMessage(*cm)
	//1 time check
	t := time.Unix(0, consensusMessage.SequenceID)
	elapsed := time.Since(t)

	if math.Abs(elapsed.Minutes()) > 5.0 {
		logger_pbftservice.Errorln("time over (5min)")
		return
	}

	logger_pbftservice.Infoln("Time check OK")

	//2 consensus id check
	cs.Lock()
	defer cs.Unlock()

	logger_pbftservice.Infoln(consensusMessage.SenderID)

	consensusID := consensusMessage.ConsensusID
	consensusState, ok := cs.consensusStates[consensusID]

	logger_pbftservice.Infoln("Message type is:",consensusMessage.MsgType)

	if !ok{
		//consensus state생성
		//prepremessage인 경우에만 block과 view, stage를 세팅
		//var newConsensusState *domain.ConsensusState
		consensusStateBuilder := domain.NewConsensusStateBuilder()

		if consensusMessage.MsgType == domain.PreprepareMsg{
			logger_pbftservice.Infoln("prepreparedmsg")
			consensusState = consensusStateBuilder.
				ConsensusID(consensusMessage.ConsensusID).
				CurrentStage(domain.PrePrepared).
				View(&consensusMessage.View).
				Block(consensusMessage.Block).
				EndConsensusHandler(cs.EndConsensusState).
				Period(300).Build()

		}else{
			logger_pbftservice.Infoln("not prepreparedmsg")
			consensusState = consensusStateBuilder.
				ConsensusID(consensusMessage.ConsensusID).
				EndConsensusHandler(cs.EndConsensusState).
				Period(300).Build()
		}

		cs.consensusStates[consensusState.ID] = consensusState
	}

	logger_pbftservice.Infoln("Add message to consensusState")
	consensusState.AddMessage(*consensusMessage)

	logger_pbftservice.Infoln("Current Stage is",consensusState.CurrentStage)

	if consensusState.CurrentStage == domain.PrePrepared{
		logger_pbftservice.Infoln("my id", cs.identity.PeerID)
		sequenceID := time.Now().UnixNano()
		logger_pbftservice.Infoln("block", consensusState.Block)
		preprepareConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,*consensusState.View,sequenceID,consensusState.Block,cs.identity.PeerID,domain.PrepareMsg)
		consensusState.CurrentStage = domain.Prepared
		cs.broadcastMessage(preprepareConsensusMessage)
	}

	//1. prepare stage && prepare message가 전체의 2/3이상 -> commitMsg전파
	if consensusState.CurrentStage == domain.Prepared && consensusState.PrepareReady(){
		sequenceID := time.Now().UnixNano()
		commitConsensusMessage := domain.NewConsesnsusMessage(consensusState.ID,*consensusState.View,sequenceID,consensusState.Block,cs.identity.PeerID,domain.CommitMsg)
		consensusState.CurrentStage = domain.Committed
		cs.broadcastMessage(commitConsensusMessage)
	}

	//2. commit state && commit message가 전체의 2/3이상 -> 블록저장
	if consensusState.CurrentStage == domain.Committed && consensusState.CommitReady(){
		logger_pbftservice.Infoln("ConsensusState is End")

		//block 저장
		//todo block에 저장
		flag, err := cs.blockService.VerifyBlock(consensusState.Block)

		if err != nil{
			common.Log.Error("Verify block error:",err.Error())
		}

		if flag{
			common.Log.Debugln("Add block")
			_, err := cs.blockService.AddBlock(consensusState.Block)

			if err !=nil{
				common.Log.Error(err.Error())
			}
			cs.transactionService.DeleteTransactions(consensusState.Block.Transactions)
		}

		cs.EndConsensusState(consensusState)
		return
	}
}

func (cs *PBFTConsensusService) EndConsensusState(consensusState *domain.ConsensusState){

	cs.Lock()
	defer cs.Unlock()

	cs.consensusStates[consensusState.ID].End()
	delete(cs.consensusStates,consensusState.ID)
}

//tested
func (cs *PBFTConsensusService) broadcastMessage(consensusMsg domain.ConsensusMessage){

	logger_pbftservice.Infoln("broadcast Message")
	peerIDList := consensusMsg.View.PeerID

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_ConsensusMessage{
		ConsensusMessage: domain.ToConsensusProtoMessage(consensusMsg),
	}

	for _, peerID := range peerIDList{
		logger_pbftservice.Infoln("sending...",peerID)
		cs.comm.SendStream(message,nil,nil,peerID)
	}
}