# Blockchain Component

## Intro

Blockchain Component는 blockchain 동기화, block 추가, 조회를 담당한다.




## API

* **Synchronize**
  * 블록을 동기화시킨다. 여기서 동기화체크(syncedCheck), 구축(constructChain), 재구축과정(addBlocksFromBlockPool)을 포함한다.



## Event

* **NodeUpdateEvent**
  * Reliable Node의 후보가 바뀔 때 바뀐 Reliable Node 후보 정보 event를 발생시킨다.



## Command Handler

* **HandleNodeUpdateCommand**
  * Reliable Node의 후보를 저장한다.



## Event Handler

* **HandleNodeUpdatedEvent**
  * Reliable Node의 후보를 변경한다.



##  Message Handler Protocol

* **BlockRequestProtocol**
* **BlockResponseProtocol**



## Command Service

* **SendSyncUpdateCommand**
  * 블록체인 컴포넌트의 동기화 상태을 알리고 동기화가 완료되었을 경우 블록체인의 마지막 블록을 보낸다.



## Message Service

* **RequestBlock**
  * Reliable Node에게 block을 요청한다.
* **ResponseBlock**
  * 요청을 보낸 Node에게 block을 보낸다.

## Block Synchronize 

1. <u>동기화(Synchronize)</u>는 특정 노드의 블록 체인을 네트워크 내의 합의된 블록 체인과 동일하게 만드는 과정을 의미한다. 즉 <u>동기화(Synchronize)</u> 과정을 통해 특정 노드는 모든 블록에 대하여 대표값(Seal), 이전 블록의 대표값(PrevSeal), 트랜잭션 모음(TxList), 트랜잭션 대표값(TxSeal), 블록 생성 시각(TimeStamp), 생성자(Creator), 블록 체인의 길이(Height) 등의 블록 체인과 관련된 모든 정보들을 다른 노드들의 것과 동일화한다.
2. <u>동기화(Synchronize)</u>는 **확인(Check)**, **구축(Construct)**, **재구축(PostConstruct)** 의 과정을 거친다.
3. **확인(Check)** 은 특정 노드의 블록 체인이 동기화가 필요한 상태인지를 점검한다. **확인(Check)** 의 과정은 임의의 노드에게 Blockchain 길이와 lastSeal을 받아와서 자신의 블록 체인 정보가 같은 지 비교하여, 동기화가 필요한 상태인지 점검한다(SyncedCheck). 이미 동기화가 완료된 상태라면, <u>동기화(Synchronize)</u> 과정을 중단한다. 그렇지 않을 경우, **구축(Construct)** 을 수행한다.
4. **구축(Construct)** 은 **확인(Check)** 과정에서 확정했던 하나의 노드에게 블록 정보를 요청(RequestBlock)하여, 응답(ResponseBlock) 받고, 응답받은 블록을 블록 체인에 추가(BlockAddedEvent)하는 과정을 순차적으로 반복함으로써 수행된다.
5. 블록 요청(RequestBlock)은 특정 노드의 블록 체인 길이(Height)를 활용해, 신뢰할 수 있는 하나의 노드에 블록을 요청함으로써 수행된다. 특정 노드가 새로 참여하는 노드일 경우 신뢰할 수 있는 노드의 블록 체인 내 최초 블록부터 마지막 블록까지 요청하고, 기존에 참여중이던 노드일 경우 보유 중인 블록 체인 내 마지막 블록의 다음 블록부터 신뢰할 수 있는 노드의 블록 체인 내 마지막 블록까지 요청한다.
6. 블록 응답(ResponseBlock)은 전달 받은 블록의 길이(Heigt)를 활용하여 블록을 조회하여(GetBlockByHeight), 요청한 노드에게 응답함으로써 수행된다.
7. 블록 추가(BlockAddedEvent)는 블록 요청(RequestBlock), 블록 응답(ResponseBlock)이 완료된 후 수행된다. 모든 블록이 추가되면 **구축(Constrcut)**이 완료된다.
8. 특정 노드는 **구축(Construct)** 의 진행 중에 새롭게 합의되는 블록을 블록 임시 저장소(BlockPool)에 보관한다. **구축(Construct)** 이 완료되고 나면, 블록 임시 저장소에 블록이 보관되어 있는 지 확인한다(PoolCheck). 보관중인 블록이 있다면, **재구축(PostConstruct)** 을 수행한다.
9. **재구축(PostConstruct)** 은 이미 **구축(Construct)** 된 블록 체인에 블록 임시 저장소(BlockFool)에 보관중인 블록들을 부수적으로 추가(BlockAddedEvent)하는 것을 의미한다. **재구축(PostConstrcut)** 을 수행하고 나면, <u>동기화(Synchronize)</u> 과정이 모두 완료된다.

## Implementation Details

[IMPLEMENTATION-DETAILS-KR.md](PROJECT-IMPLEMENTATION-DETAILS-KR.md)

### Author

[@junk-sound](https://github.com/junk-sound), [@zeroFruit](https://github.com/zeroFruit)
