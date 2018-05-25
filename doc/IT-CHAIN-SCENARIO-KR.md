---
본 문서는 peer 개발 시 소통을 위해 임의로 작성되었으며
전체 개발의 진행방향과 무관하게 소통을 위한 문서이기에 혼동이 없으셨으면 합니다.
-frontalnh(namhoon)
---

# requirement
rabbitmq 설치 및 서버 가동

# 최초의 peer에서 블록체인 초기화
## 부트 노드 환경 설정
**자신이 부트노드인 경우**
peer/init.go 에서 init() 함수 호출시 부트노드 및 나의 노드에 대한 정보를 설정하고 해당 노드가 부트노드인 경우 자신을 리더로 선정함.


**자신이 부트 노드가 아닌 경우**
별도의 authenticate 절차 거침

## grpc 서버 개통 및 p2p 네트워크 연결
gateway 컴포넌트에서 Start() 호출하여 p2p 네트워크에 연결한다.
## amqp 서버 연결
```
messaging/rabbitmq
MessageQueue.start()
```
## Gateway 연결
gateway.init()



# 두번째 피어 등록
자신이 부트노드가 아니므로
1. 피어 생성
2. 다른 피어에게 알림
3. 기존 피어목록 생신
4. 내 피어 생성

---
아래 내용은 소통을 위해 생각의 흐름대로 임의로 작성하였으며
추후 개선될 예정입니다.
---


# 트랜잭션 발생
1. 사용자가 tx 생성 요청
2. txpool 에 저장
3. txpool에서 consume
4. pool에 등록
5. api 호출
createtx
6. amqp 등록
7. create event 객체생성
8. pool에서 tx 생성
9. leveldb 저장
# 블록 생성
# 합의과정
# 블록 저장
# 블록 전파
