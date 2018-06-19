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

* **SendSyncStartCommand**
  * 블록체인 컴포넌트의 동기화 시작을 알린다.
* **SendSyncFinishCommand**
  * 블록체인 컴포넌트의 동기화가 완료되었다는 것을 알리고, 블록 체인 내 마지막 블록을 보낸다.



## Message Service

* **RequestBlock**
  * Reliable Node에게 block을 요청한다.
* **ResponseBlock**
  * 요청을 보낸 Node에게 block을 보낸다.
