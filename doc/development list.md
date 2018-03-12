# IT-CHAIN 개발 리스트



## BlockChain - 블록체인 관련 구조체 및 함수 구현 

- Blockchain 관련 함수

- block 생성, 검증, merkle tree

- transaction 생성, 검증

  ​

## P2P - 다른 peer(node)와의 여러 interaction

- Consensus 합의 알고리즘

- Gossip protocol로 peer 네트워크 유지 및 관리

  ​

## DB - Blockchain, peer_list, smart-contract 저장 및 관리  

- level DB로 blockchain 저장, 복구, 관리

- world state DB(최신 현재 상태의 정보) 관리

  ​

## Authentication - 디지털 서명 검증 과 peer 인증 관리

- Block,transaction 서명 및 검증

- 내 private key, public key 관리

- peer들의 인증 관리 

  ​

## VM - smart-contract 실행 환경 구성 및 실행

- 가상 환경에서 Docker를 이용한 transaction 실행 

  ​

## Service - grpc 서비스 로직 구현

-  peer 찾기 api

- 트랜잭션 생성 api

- blockchain 통계 정보 api

- consensus api

- 트랜젹선 검증 api

  ...

## Smart contract repository - 스마트 컨트랙 저장 및 관리

- 깃과 연동하여 smart contract 등록 및 관리 인증 처리

