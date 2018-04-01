# it-chain


<p align="center"><img src="./images/logo.png" width="300px" height="200px"></p>

## Overview

Generalized Private Chain For All 

it-chain makes it easy for anyone to have their own block chain platform

Use Smart-Contract to create your own Dapp



## Requirements

- Go-lang 1.9



## Config

location : ./conf/config.yaml

```
database:
  type: leveldb
  leveldb:
    defaultPath: .it-chain/leveldb

key:
  type: RSA
  defaultPath: .it-chain/

ledger:
  defaultPath: .it-chain/ledger

txDatabase:
  defaultPath: .it-chain/tx

smartContract:
  defaultPath: .it-chain/.constracts
  githubID: {GithubID}
  githubPassword: {GithubPassword}
  githubAccessToken: {GithubAccessToken}
  TmpDir: /tmp
  WorldStateDB_OnDoker: /go/src/worldstatedb

batchTimer:
  pushPeerTable: 5
  sendTransaction: 5

//The IP of the BootNode becomes the leader.
bootNode:
  ip: 127.0.0.1:4444

//The node is deployed to the IP of the node. 
node:
  ip: 127.0.0.1:4444

consensus:
  batchTime: 3
  maxTransactions: 100

webhook:
  port: 44444
```



## Usage

```
go run main.go
```


## Blockchain Key Concept
The key concept of blockchain can be found in the KEYCONEPT. <br>
[KEYCONCEPT](doc/KEYCONCEPT.md)

## Logical Architecture of `it-chain`
Overview of the logical architecture of `it-chain` can be found in the following link. At this point, only the Korean version is provided.<br>
[LOGICAL ARCHITECTURE KR](doc/LOGICAL-ARCHITECTURE-KR.md)

## Implementation Details
Core implementation decisions can be found in the Project Implementation Details. <br>
[PROJECT IMPLEMENTATION DETAILS](doc/PROJECT-IMPLEMENTATION-DETAILS.md)

## Contribution
Contribution Guide <br>
[CONTRIBUTION](CONTRIBUTION.md)

## License

It-Chain Project source code files are made available under the Apache License, Version 2.0 (Apache-2.0), located in the [LICENSE](LICENSE) file.

## Designed by
@Hyemin choi<br>
@Jieun Oh<br>
@Jongmo Moon<br>
