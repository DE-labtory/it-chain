# 블록체인 동기화

1. 동기화(Synchronize)는 특정 노드의 블록 체인을 네트워크 내의 합의된 블록 체인과 동일하게 만드는 과정을 의미한다. 즉 동기화 과정을 통해 특정 노드는 모든 블록에 대하여 대표값(Seal), 이전 블록의 대표값(PrevSeal), 트랜잭션 모음(TxList), 트랜잭션 대표값(TxSeal), 블록 생성 시각(TimeStamp), 생성자(Creator), 블록 체인의 길이(Height) 등의 블록 체인과 관련된 모든 정보들을 다른 노드들의 것과 동일화한다.

1. 동기화는 확인(Check), 구축(Construct), 재구축(PostConstruct)의 과정을 거친다.
2. 확인(Check)은 특정 노드의 블록 체인이 동기화가 필요한 상태인지를 점검한다. 확인(Check)의 과정은 (1) '합의된 블록 체인 정보 가져오기', (2) '상태 점검' 을 거친다. (1) '합의된 블록 체인 정보 가져오기' 는  먼저 P2P 네트워크 내의 다양한 노드에게 각자 보유중인 블록 체인의 정보를 요청한다(BlockChainLoadCommand). 이 때 블록 체인 전부를 요청하는 것이 아니라, 확인(Check)에 필요한 정보만을 요청한다. 여러 블록 체인 정보를 불러오고 나면(BlockChainLoadedEvent), 불러온 정보들을 활용하여 신뢰할 수 있는 하나의 노드 및 블록 체인을 확정한다(ReliableNodeConfirm, BlockChainConfirm). 이 후 (2) '상태 점검' 은 특정 노드의 블록 체인 정보와 합의된 블록 체인 정보가 같은 지 비교하여, 동기화가 필요한 상태인지 점검한다(StateCheck). 이 때, 빠르게 비교할 수 있는 순서대로 블록 체인의 길이(Height), 마지막 블록의 대표값(Seal), 모든 블록의 대표값(Seal) 등을 순차적으로 비교한다. 이미 동기화가 완료된 상태라면(isSynchronized), 동기화(Synchronize) 과정을 중단한다. 그렇지 않을 경우, 구축(Construct)을 수행한다.
3. 구축(Construct)은 합의된 블록 체인의 모든 블록을 순회하며 블록 요청(BlockRetrieveCommand), 블록 응답(BlockRetrievedEvent)의 과정을 반복하고, 응답받은 모든 블록들을 한꺼번에 추가(BlockAdd)함으로써 수행된다.
4. 블록 요청(BlockRetrieveCommand)은 확인(Check) 과정에서 확정했던 신뢰할 수 있는 하나의 노드에 블록의 길이(Height)를 기준으로 블록을 요청함(GetBlockByHeight)으로써 수행된다. 특정 노드가 새로 참여하는 노드일 경우(isNewnode) 합의된 블록 체인 내 최초 블록부터 마지막 블록까지 요청하고, 기존에 참여중이던 노드일 경우 보유 중인 블록 체인 내 마지막 블록의 다음 블록부터 합의된 블록 체인 내 마지막 블록까지 요청한다.
5. 블록 응답(BlockRetrievedEvent)은 블록 요청(BlockRetrieveCommand)에 따른 결과이다.
6. 블록 요청(BlockRetrieveCommand)과 블록 응답(BlockRetrievedEvent)의 반복으로 추가해야 할 블록들을 모두 불러오고 나면(isDownloadFinished), 특정 노드의 블록 체인에 블록들을 모두 추가한다(BlockAdd). 모든 블록을 추가하고 나면 구축(Construct)이 완료된다.
7. 특정 노드는 구축(Construct)의 진행 중에 새롭게 합의되는 블록을 블록 임시 저장소(BefferFool)에 보관한다. 구축(Construct)이 완료되고 나면, 블록 체인의 완결성을 보장하기 위해 블록 임시 저장소에 블록이 보관되어 있는 지 확인한다(BufferCheck). 보관중인 블록이 있다면, 재구축(PostConstruct)을 수행한다.
8. 재구축(PostConstruct)은 이미 구축된 블록 체인에 블록 임시 저장소(BufferFool)에 보관중인 블록들을 부수적으로 추가(BlockAdd)하는 것을 의미한다. 재구축(PostConstrcut)을 수행하고 나면, 동기화(Synchronize) 과정이 모두 완료된다.