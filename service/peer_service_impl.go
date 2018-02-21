package service

import (
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"it-chain/domain"
	"strconv"
	"time"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"it-chain/network/comm/msg"
	"it-chain/common"
	"github.com/golang/protobuf/proto"
	"it-chain/network/comm/conn"
)

type PeerServiceImpl struct {
	peerTable *domain.PeerTable
	comm      comm.ConnectionManager
}

func NewPeerServiceImpl(peerTable *domain.PeerTable,comm comm.ConnectionManager) PeerService{

	peerService := &PeerServiceImpl{
		peerTable: peerTable,
		comm: comm,
	}

	i, _ := strconv.Atoi(viper.GetString("batchTimer.pushPeerTable"))

	broadCastPeerTableBatcher := NewBatchService(time.Duration(i)*time.Second,peerService.BroadCastPeerTable,false)
	broadCastPeerTableBatcher.Add("push batching")
	broadCastPeerTableBatcher.Start()

	//todo 양방향 reference제거하고싶다..
	comm.SetOnConnectHandler(peerService.HandleOnConnect)
	comm.Subscribe("EstablishConnection",peerService.handleConnectionEstablish)
	comm.Subscribe("UpdatePeerTable",peerService.handleUpdatePeerTable)

	return peerService
}

func (ps *PeerServiceImpl) GetPeerTable() domain.PeerTable{
	return *ps.peerTable
}

//peer info 찾기
func (ps *PeerServiceImpl) GetPeerByPeerID(peerID string) (*domain.Peer){

	Peer := ps.peerTable.FindPeerByPeerID(peerID)
	return Peer
}

//peer info
func (ps *PeerServiceImpl) PushPeerTable(peerIDs []string){

}

//주기적으로 handle 함수가 콜 된다.
//주기적으로 peerTable의 peerlist에게 peerTable을 전송한다.
func (ps *PeerServiceImpl) BroadCastPeerTable(interface{}){
	common.Log.Println("pushing peer table")

	Peers, err := ps.peerTable.SelectRandomPeers(0.5)

	if err != nil{
		common.Log.Println("no peer exist")
		return
	}

	ps.peerTable.IncrementHeartBeat()

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_PeerTable{
		PeerTable: domain.ToProtoPeerTable(*ps.peerTable),
	}

	if err !=nil{
		common.Log.Println("fail to serialize message")
	}

	errorCallBack := func(onError error) {
		common.Log.Println("fail to send message error:", onError.Error())
	}

	for _,Peer := range Peers{
		ps.comm.SendStream(message,nil,errorCallBack, Peer.PeerID)
	}
}

func (ps *PeerServiceImpl) UpdatePeerTable(peerTable domain.PeerTable){

	ps.peerTable.Lock()
	defer ps.peerTable.Unlock()

	for id, Peer := range peerTable.PeerMap{
		peer,ok := ps.peerTable.PeerMap[id]
		if ok{
			peer.Update(Peer)
		}else{
			ps.AddPeer(Peer)
		}
	}

	ps.peerTable.UpdateTimeStamp()
}

func (ps *PeerServiceImpl) AddPeer(Peer *domain.Peer){

	if Peer.PeerID == ""{
		common.Log.Error("failed to connect with", Peer)
		return
	}

	if Peer.GetEndPoint() == ""{
		common.Log.Error("failed to connect with", Peer)
		return
	}

	err := ps.comm.CreateStreamClientConn(Peer.PeerID,Peer.GetEndPoint())

	if err != nil{
		common.Log.Error("failed to connect with", Peer)
		return
	}

	ps.peerTable.AddPeer(Peer)
}

func (ps *PeerServiceImpl) RequestPeer(ip string) (*domain.Peer ,error){

	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		common.Log.Println("can not connect: %v", err)
		return &domain.Peer{},err
	}

	defer conn.Close()
	c := pb.NewPeerServiceClient(conn)

	peer, err := c.GetPeer(context.Background(), &pb.Empty{})

	if err != nil {
		common.Log.Println("fail to get peerInfo: %v", err.Error())
		return &domain.Peer{},err
	}

	return domain.FromProtoPeer(*peer), nil
}

func (ps *PeerServiceImpl) GetLeader() *domain.Peer{
	return ps.peerTable.Leader
}

func (ps *PeerServiceImpl) HandleOnConnect(conn conn.Connection, pp pb.Peer){
	peer := domain.FromProtoPeer(pp)
	ps.peerTable.AddPeer(peer)
}

func (ps *PeerServiceImpl) handleConnectionEstablish(message msg.OutterMessage){

	if establishMsg := message.Message.GetConnectionEstablish(); establishMsg != nil{
		common.Log.Println("Handling connection establish")
		respondMessage := &pb.StreamMessage{}
		respondMessage.Content = &pb.StreamMessage_Peer{
			Peer: domain.ToProtoPeer(ps.peerTable.GetMyInfo()),
		}

		payload, err := proto.Marshal(respondMessage)

		if err !=nil{
			common.Log.Println("Marshal error:", err.Error())
		}

		respondEnv := &pb.Envelope{}
		respondEnv.Payload = payload

		var errCallBack = func (err error){
			common.Log.Println("Respond error:", err.Error())
		}

		message.Respond(respondEnv,nil,errCallBack)
	}

	return
}

func (ps *PeerServiceImpl) handleUpdatePeerTable(message msg.OutterMessage){

	if peerTableMsg := message.Message.GetPeerTable(); peerTableMsg != nil{
		common.Log.Println("Handling peertable update message")

		peerTable := domain.FromProtoPeerTable(*peerTableMsg)
		ps.UpdatePeerTable(*peerTable)

		common.Log.Println("PeerTable:",ps.peerTable)
	}

	return
}

func (ps *PeerServiceImpl) SetLeader(peer *domain.Peer){
	ps.peerTable.Leader = peer
}