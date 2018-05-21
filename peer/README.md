# 서론
Peer 서비스는 피어노드의 생성, 삭제, 열람, 리더 노드의 선정 및 변경 등의 역할을 수행합니다.

# API

## LeaderSelection Structure
peerTableService, messageProducer, peerRepository, myInfo 속성을 가지며 각각은 다음의 역할을 수행한다.


## Functions of LeaderSelection Structure
### NewLeaderSelectionApi()
리더 선출 관련 api 인 leaderSelectionApi 구조체를 생성한다.

**input**
`repo repository.Peer`
`messageProducer service.MessageProducer`
`myInfo *model.Peer`

**output**
`leaderSelectionApi`
`nil`


### RequestChangeLeader()
구현안됨

### RequestLeaderInfoTo()
특정 피어를 입력으로 전달하여 해당 피어에게 리더의 정보를 가르쳐 준다.

**input**
`peer model.Peer`


**output**
error

### changeLeader()

peerTableService의 setLeader를 호출한다.

**input**
`peer *model.Peer`

**output**
error if exist


# domain

## model
model의 구조체를 정의하고 value 체크가 가능한 수준의 validate 함수 및 각종 변환함수들이 구현됨.

## repository
infra repository에 대한 interface만 구현

## service
peer 와 관련된 다양한 기능 수행
실제 구현까지 이루어져야 함.

# infra
## messaging
## repository
db 통신을 하는 실제 구현이 이루어짐.

현재 leveldb에 접근하는 용도로 사용


---
