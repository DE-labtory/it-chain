# Config

Config 파일을 생성하고 Config 값을 불러오는 것에 관련한 패키지.

Config는 세개의 주요 기능으로 구성되어있다.

- Config에 담길 데이터를 정의하는 model
- Config파일을 생성하는 generator
- Config파일을 읽는 Configuration

## Model
Config 에 들어갈 데이터 모델들이 정의되어있다.

정의되어 있는 모델은 다음과 같다.
- Authentification
- Blockchain
- Consensus
- GrpcGateway
- ICode
- Peer
- Txppol
- Common

각 모델 안에 필요한 값들이 Struct 형태로 저장되어있다.

## Generator
config 모델을 생성/변경 후 **Config.yaml** 파일을 생성/변경 하기 위한 모듈이다.

config 모델을 생성/변경 후 개발 환경에 맞는 스크립트 파일을 실행해주면 된다

Linux/Unix(Mac OS)

     user:~/go/src/github.com/DE-labtory/engine/conf$ sh generate_conf.sh

Windows

    c:/~~~/go/src/github.com/DE-labtory/engine/conf> generate_conf.bat
    또는
    generate_conf.bat 을 클릭하여 실행


## Configuration (config usage)
다른 패키지에서 conf를 임포트하여 사용할 수 있다.

GetConfiguration 함수를 호출하여 config 데이터를 불러올 수 있다.

    import "github.com/DE-labtory/engine/conf"

    func foo(){
        config = conf.GetConfiguration()
        defaultPath = config.Txpool.RepositoryPath;
        ...
        ...
    }

## How to add/change config data
- 변경 방법

  - ./conf/model 에 있는 값들을 변경하고 generator를 빌드/실행하면 된다.

- 추가 방법
  - ./conf/model 에 존재하는 모델에 추가할경우
  <br>모델에 구조체/필드 를 추가하고 generator를 빌드/실행하면 된다.

  - ./conf/model 에 존재하지 않고 모델을 새로 추가할 경우
  <br>./conf/model에 모델 go 파일을 생성한 뒤 객체를 반환하는 NewModel함수를 만든다.
  <br>그 후 ./conf/configuration.go 에 Configuration Struct에 해당 모델을 추가하고
  <br>./conf/main/conf_file_generator에 Configuration Struct를 할당하는부분에 NewModel 함수를 호출해준다.
