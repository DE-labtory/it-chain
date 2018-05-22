# 최초의 peer에서 블록체인 초기화

1. peer 생성
2. AMQP 서버 개통
```
messaging/rabbitmq
MessageQueue.start()
```
3. Gateway 연결
gateway.init()

4.

# 두번째 피어 등록
1. 피어 생성
2. 다른 피어에게 알림
3. 기존 피어목록 생신
4. 내 피어 생성
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
