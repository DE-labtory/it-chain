package p2p

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/common"
)

var NodeNotFoundErr = errors.New("NodeNotFound")
var DuplicateNodeErr = errors.New("NodeAlreadyExist")
var NodeIdEmptyErr = errors.New("NodeIdEmpty")

// 피어 구조체를 선언한다.
type Node struct {
	IpAddress string
	Id        NodeId
}

// 해당 피어의 ip와 peerId로 새로운 피어를 생성한다.
func NewPeer(ipAddress string, id NodeId) *Node {
	return &Node{
		IpAddress: ipAddress,
		Id:        id,
	}
}

// p2p 구조체를 json 으로 인코딩한다.
func (n Node) Serialize() ([]byte, error) {
	return common.Serialize(n)
}

// 입력받은 p2p 구조체에 해당 json 인코딩 바이트 배열을 deserialize 해서 저장한다.
func Deserialize(b []byte, peer *Node) error {
	err := json.Unmarshal(b, peer)

	if err != nil {
		return err
	}

	return nil
}

// peerId 선언
type NodeId string

// conver peerId to String
func (nodeId NodeId) ToString() string {
	return string(nodeId)
}

// p2p repository 인터페이스를 정의한다.
// NodeRepository 가 아니라 NodeRepository 로 정의하는 것이 맞을것 같습니다. - 남훈
type NodeRepository interface {
	Save(node Node) error
	Remove(id NodeId) error
	FindById(id NodeId) (*Node, error)
	FindAll() ([]*Node, error)
}
