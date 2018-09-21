해당 문서에서 다루는 목록은 다음과 같다

- IVM(링크)
- ICode(링크)
- ConsumeEvent/Command(링크)
- PublishEvent/Command(링크)
- FutureWork(링크)

# IVM (Icode Virtual Machine)

Icode Virtual Machine(IVM) 은 it-chain 내에서 ICode[아래로링크]를 실행시키고 관리하기 위한 컴포넌트이다.

It-chain 프로젝트 내 Tesseract 라이브러리를 통하여 ICode들의 실행 및 동작을 관리한다.

아래 그림은 IVM의 관계 구조이다.

[그림 - IVM-teseeract - (sdk - icode)여러개]

IVM은 tesseract를 이용하여 여러개의 ICode를 관리하며, 이때 각 ICode와 통신은 ICode에서 사용한 SDK[링크]와 통신한다.

각 ICode들은 Tesseract를 통해 Docker 가상환경에서 관리되고 이를 Container라고 부른다.



ICode를 IVM위에서 실행시키기 위해서는 먼저 IVM에 ICode를 Deploy하는 과정이 필요하다.

아래 그림은 IVM에 ICode를 Deploy 하는 순서도이다

[그림 - ICode Deploy 과정 sequence diagram]

현재 it-chain은 git protocol을 통해 deploy 하는것을 지원한다.



IVM위에 올라간 ICode는 Transaction에 의해 실행된다.

Transaction은 블록에 포함되는데 블록이 합의[consesnsus doc 링크]된 후 블록 내 모든 Transaction을 일괄 처리한다.

블록이 합의된 후 블록 내 트랜젝션이 실행되는 순서도는 다음과 같다.

[그림 - Block에서 Transaction들 ICode 실행하는 순서도]

블록 내에 존재하는 트랜젝션들을 timestamp 순서에 맞게 실행시킨다.

하나의 트랜젝션 안에는 ICode정보와 Function, Argument 정보가 들어있다.

ICode정보를 이용하여 IVM이 해당 Container를 찾아 Argument와 Function정보를 전달하고 결과를 요청한다.



## ICode

icode 는 it-chain 위에서 실행되는 스마트 컨트랙트이다.

icode를 작성하기위한 SDK가 제공되며 다양한 언어를 지원한다.

SDK에 대한 자세한 정보는 <여기>[SDK 페이지로 링크] 를 통해 확인 할 수 있다.

ICode의 Function Type은 Invoke와 Query로 나뉜다

Invoke Function의 경우 ~~~~

[Invoke Function type 설명]

Query Function의 경우 ~~~

[Query Function type 설명]



## Consume Event/Command

아이코드가 Consume 하는 이벤트/커맨드에 대해 다룸.



## Publish Event/Command

아이코드가 Publish 하는 이벤트/커맨드에 대해 다룸.



## Future Work

현재 stable하게 동작하는 IVM - tesseract - SDK에서 추가할 기능은 다음과 같다.

- 공통 World State DB[아래로 링크]
- Read / Write Permission[아래로 링크]
- ICode version control[아래로 링크]
- Result Set Validation[아래로 링크]



### 공통 World State DB

현재 IVM 위에서 동작하는 icode들은 독립적인 데이터베이스를 가지고 동작한다.

이를 tesseract가 관리하는 공동데이터베이스로 바꾸고 한곳에서 world-state database를 구성하고자한다.

이는 아래 나올 Read / Write Permission 및 Result Set Validation 과도 관련이 있다.



### Read / Write Permission

 향후 IVM에 world-state database가 구성이된다면 서로 다른 Icode들이 만들어낸 정보를 접근 할 수 있게된다.

이때 무분별한 Read/Write를 막고자 Icode 에서 생성한 정보의 Read/Write 권한은 해당 Icode에 명시된 Icode들만 허용된다.



### ICode Version Control

스마트컨트랙트 즉  ICode는 현재 업데이트가 되면 다른 ICode로인식되어 새로 Deploy 하여야 한다. ICode의 버전 업데이트 방식을 지원하는것이 하나의 목표이다.



### Result Set Validation

합의 과정에서 블록 내의 트랜젝션이 변경하는 값들을 비교하고 합의하고자 한다.

블록 내의 트랜젝션들은 여러 ICode들을 실행하게 되는데 이에 따른 State들의 변경사항을 기록해 두어 블록 합의에 담는 것을 목표로 한다.

이는 다양한 언어를 지원하는 it-chain SDK의 특성상 ICode 의 실행결과값이 항상 같다는 것을 보장해주지 못하는 것을 막고,

모든 노드가 같은 State를 가지도록 보장해준다.