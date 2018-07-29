# Consensus

합의(Consensus)는 engine의 P2P Network의 구성원들 사이에서 수행되며 리더노드에 의해 시작된다. Consensus Module은 리더에 의해 시작된 **Block의 저장 순서에 대한 합의** 만을 수행한다.

블록에 포함된 내용이 이미 잘 검증된 트랜잭션들이고, 각 트랜잭션의 결과가 deterministic 하다면, Block의 저장 순서만 정확히 합의할 경우 모든 노드가 언제나 동일한 결과물(world state)을 갖게 된다.

## pbft consensus algorithm in it-chain engine

In it-chain core engine, the consensus is made by PBFT algorithm.

![](../doc/images/consensus-PBFT.png)

the following is a brief explanation of the procedure of PBFT consensus algorithm

1. the client requests to every node in the network
2. if the leader of the network receives the request of a client, it broadcasts preprepare messages to the p2p network
3. every node in the network which received the message broadcasts prepare message to the network
4. a node which received more than 2/3 of prepare messages broadcasts commit message to the network
5. consensus complete

## Basic Concept of implementation

1. leader receives the request from the client
2. leader create the consensus that has information about representatives 3. broadcasts preprepare messages to representatives
4. the node who receives preprepare message construct the consensus by given info and broadcast prepare message to the network
5. a node which received more than 2/3 of prepare messages broadcasts commit message to the network. until the node receives all representatives prepare message, the node saves the prepare message in prepare message pool
6. every node who receives commit message saves the messages until the received block is confirmed

### Parliament

`Parliament` is the group of nodes which participate in consensus procedure. every node in parliament is called `representative` in the sense of being a voter in consensus.

### Message pool

**prepare message pool**
Before broadcast commit message, every node should determine the message that includes the most reliable block that most of the nodes agree with.
so that the `prepare message` should be saved in `prepare message pool` until the node receives every other node prepare message 

**commit message pool**
until the block is confirmed, the commit message should be saved in commit message pool

## procedure detail

1. leader receives request from client
2. leader create the consensus that has information about representativesbroadcasts preprepare messages
3. the node who receives 

## kind of pbft consensus messages

The consensus between nodes is made by sending and receiving certain kind of consensus messages

The following is the kind of message used between nodes.

### Preprepare message

The message sent by leader of network to inform other nodes in network 

### prepare message

The message sent to every other node from every node in network.

every node in the network of n nodes receives the number of n-1 prepare messages so that every node in the network can determine which is the proper message which is chosen by most of the nodes

### commit message

Every node in the network determine the most reliable message which is chosen by most of the nodes and broadcast commit message to every node in the network.

## Behavior at event occurrence

### LeaderChanged

### MemberJoined

### MemberRemoved

### BlockSaved

## Event Publishing

- Consensus 관련
  - ConsensusPrePrepared (Preprepare Msg를 보냈을 때)
  - ConsensusPrepared (Prepare Msg를 보냈을 때)
  - ConsensusCommitted (Commit Msg를 보냈을 때)
  - ConsensusFinished (블록을 저장했을 때 IDLE 상태가 된다)

- ConsensusMessageAddedEvent (3 types)
  - PreprepareMsg
  - PrepareMsg
  - CommitMsg



## Consume Command

- PreprepareMsg 받음
- PrepareMsg 받음
- CommitMsg 받음

## Publish Command

- CreateBlockCommand
- SendConsensusMsgCommand
  - PreprepareMsg (leader)
  - PrepareMsg
  - CommitMsg
  
  

## Author
[@hihiboss](https://github.com/hihiboss)
[@frontalnh](https://github.com/frontalnh)
