# Consensus

합의(Consensus)는 engine의 P2P Network의 구성원들 사이에서 수행되며 리더노드에 의해 시작된다. Consensus Module은 리더에 의해 시작된 **Block의 저장 순서에 대한 합의** 만을 수행한다.

블록에 포함된 내용이 이미 잘 검증된 트랜잭션들이고, 각 트랜잭션의 결과가 deterministic 하다면, Block의 저장 순서만 정확히 합의할 경우 모든 노드가 언제나 동일한 결과물(world state)을 갖게 된다.

## PBFT consensus algorithm

In it-chain core engine, the consensus is made by PBFT algorithm.

![](../doc/images/consensus-PBFT.png)

The following is a brief explanation of the procedure of PBFT consensus algorithm

1. the client requests to every node in the network
2. if the leader of the network receives the request of a client, it broadcasts preprepare messages to the p2p network
3. every node in the network which received the message broadcasts prepare message to the network
4. a node which received more than 2/3 of prepare messages broadcasts commit message to the network
5. consensus complete

**In It-chain, only the leader can propose and create the block. So, there is no client, no request, and no response.**

## Basic Concept of implementation

1. Leader's consensus component receives the request from the blockchain component.
2. Leader creates the consensus that has information about representatives and proposed block.
3. Broadcasts pre-prepare messages to every representatives.
4. Each representative who received the pre-prepare message constructs the consensus by given info. Then, broadcasts prepare messages to the network.
5. Each representative who received more than "1/3 of representatives + 1" preprare messages broadcasts commit messages to the network. Until the representative receives all prepare messages, saves them in the prepare message pool.
6. Each representative who received more than "1/3 of representatives + 1" commit messages publishes the block confirm event and removes the consensus. Until the proposed block is confirmed, the commit messages are saved in the commit message pool.

### Consensus State

- IDLE_STATE : A consensus is generated, but does not anything.
- PREPREPARE_STATE : The consensus (only leader) sent the pre-prepare messages.
- PREPARE_STATE : The consensus sent the prepare messages.
- COMMIT_STATE : The consensus sent the commit messages.

### Parliament

`Parliament` is the group of nodes which participate in consensus procedure. Every node in parliament is called `representative` in the sense of being a voter in consensus. `Representatives` are selected by `func Elect(parliament []MemberId) ([]*Representative, error)`.

### Message pool

**prepare message pool**
Before broadcast commit message, every representative should determine the message that includes the most reliable block that most of the representatives agree with.
So that the `prepare message` should be saved in `prepare message pool` until the representative receives every other representative prepare message 

**commit message pool**
Until the block is confirmed, the commit message should be saved in commit message pool

## Procedure detail

1. The blockchain component of the leader requests a consensus to the consensus component.
2. The leader's conensus component creates a consensus about the requested block.
3. The leader make the Pre-prepare messages which has information of leader's consensus. Then, broadcasts them to every representative.
4. Each representative who receives the leader's Pre-prepare message creates a consensus. And sends the Prepare messages to all other representatives.
5. If the number of received prepare messages is equal to or greater than "1/3 of representatives + 1", validates the block in that message. Then, the representative sends the commit messages to all other representatives.
6. If the number of received commit messages is equal to or greater than "1/3 of representatives + 1", confirms the block.
7. Removes the finished consensus.

## The kinds of PBFT consensus messages

The consensus between representatives is made by sending and receiving certain kind of consensus messages.

The following is the kind of message used between representatives.

### Preprepare message

A message sent by the leader of network to inform other representatives in parliament. 

### Prepare message

The message sent to all other representatives from each representative.

Each representative in the network of n representatives receives the number of n-1 prepare messages, so that every representative in the network can determine which is the proper message which is chosen by most of the representatives

### commit message

Every representative in the network determine the most reliable message which is chosen by most of the representatives and broadcast commit message to all other representatives in the network.

## Event & Command

### Event
- Publish
```go
// When the consensus is finished, this event will be published.
// If 'IsConfirmed' is true, the Blockchain component saves the proposed block.
type ConsensusFinished struct {
	IsConfirmed bool
}
```

### Command
- Publish
```go
// This command is used to send pre-prepare, prepare, commit messages to the network.
// The serialized messages will be included in Body.
// Distinguishing which message is received is done by Protocol.
type DeliverGrpc struct {
	midgard.CommandModel
	RecipientList []string
	Body          []byte
	Protocol      string
}
```

- Handle
```go
// This command is handled when pre-prepare, prepare, commit messages are received.
// The messages are included in Body and distinguished by Protocol.
type ReceiveGrpc struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
}
```
```go
// This command is delivered from the Blcokchain component.
// When this command is handled, the consensus is created.
// Only leader uses this command.
type StartConsensus struct {
    
}
```

## API
```go
// The leader creates a consensus and sends a pre-prepare messages to representatives.
func StartConsensus(userId consensus.MemberId, block consensus.ProposedBlock) error
```
```go
// When the pre-prepare message is delivered, the receivers construct a consensus and send a prepare message to other representatives.
func ReceivePrePrepareMsg(msg consensus.PrePrepareMsg)
```
```go
// When the prepare messages are delivered, the receivers save those messages in the prepare message pool, validate them, and send a commit message to other representatives.
func ReceivePrepareMsg(msg consensus.PrepareMsg)
```
```go
// When the commit messages are delivered, the receivers save those messages in the commit message pool, validate them, and publish the ConsensusFinished event.
func ReceiveCommitMsg(msg consensus.CommitMsg)
```



## Future Work

- World State Value validation



## Author

[@hihiboss](https://github.com/hihiboss)
[@frontalnh](https://github.com/frontalnh)
