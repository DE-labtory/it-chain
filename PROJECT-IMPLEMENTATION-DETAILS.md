# Project Implementation Details



## Overview

It describes the important implementation decisions of the it-chain. Sample code for each detailed implementation can be found in the sample folder. 



## Table of Contents

1. [BlockChain](#BlockChain)
2. [SmartContract](#SmartContract)
3. [Communication](#Communication)
4. [Crypto](#Crypto)



## BlockChain <a name="BlockChain"></a>



## SmartContract <a name="SmartContract"></a>

![smartContract-implementation-deploy](/Users/jun/go_workspace/src/it-chain/images/smartContract-implementation-deploy.png)

SmartContract is stored on git repository and is executed by the smart contract service. After testing Smart Contract in a Docker-based virtual environment, it is reflected in the actual database.

- Git

  Each Smart Contract is stored as a Git Repository.

- Docker VM

  It is a virtual environment that executes smart contracts. After the smart contract and the world state db are copied to the Docker vm, they are executed and verified virtually.

- SmartContractService

  깃과 Docker VM을 관리하는 서비스이다. 깃을 통해 스마트 컨트랙트를 푸쉬 및 클론하고 Docker VM에 world State DB와 smart contract을 copy하여 실행시킨다. 

  It is a service that manages git and Docker VM. After pushing and cloning the smart contract on the git, it copies the world state DB and smart contract to Docker VM and executes it.

  ​

#### Deploy Smart Contract Sequence Diagram

![smartContract-implementation-seq](./images/smartContract-implementation-seq.png)

The deployed user's repository is stored and managed in the Authenticated Smart Contract Repository as shown below.

| User <br />Repository <br />Path | Smart Contract <br />Repository <br />Path | Smart Contract File Path                 |
| -------------------------------- | ---------------------------------------- | ---------------------------------------- |
| A/a                              | {authenticated_git_id}/A_a               | It-chain/SmartContracts/A_a/{commit_hash} |
| B/b                              | {authenticated_git_id}/B_b               | It-chain/SmartContracts/B_b/{commit_hash} |
| C/c                              | {authenticated_git_id}/C_c               | It-chain/SmartContracts/C_c/{commit_hash} |



### Author

[@hackurity01](https://github.com/hackurity01)

## Grpc Communication <a name="Communication"></a>

<img src="./images/grpc implementation.png"></img>

Since it is complex to handle the reception and transmission of the peers' messages while maintaining the connection between the peers, the message handler is used to separate the reception and transmission and obtain the processing service using the message type.

- ConnectionManager

  The connection manager manages grpc connections with other peers in the network. Services send messages to peers through the connection manager. The reason for maintaining the connection between each peer is to make a consensus in a short time through fast block propagation.

- Grpc Client

  The grpc client is a bi-stream and can be sent or received. Both message transmission and reception are handled as go-routines, and when a message is received, it is transmitted to the MessageHandler. The client is not involved in the message content.

- MessageHandler

  The message handler receives all messages received by the grpc client, performs message validation, and forwards the message to the corresponding service according to the message type.

- Services

  Each service has a connectionManager and sends a message through the connectionManager. Services register a callback to a message of interest to the messageHandler, and process the message through a callback when the corresponding type of message is received.

### Author

[@Junbeomlee](https://github.com/junbeomlee)


## Crpyto <a name="Crypto"></a>

![crpyto-implemenation-module](./images/crpyto-implemenation-module.png)

Crypto signs and verifies the data used in the block-chain platform and manages the keys used in the process. it-chain supports rsa and ecdsa encryption method.

- KeyGenerator

  The node generates a key that matches the encryption scheme that you want to use for signing.

- KeyManager

  Stores the generated key, and loads the stored key.

- Signer

  Performs data signature.

- Verifier

  Verify the signed data.

- KeyUtil

  Perform the necessary processing tasks in the process of storing and loading the key.

<br>

### Signing process of data
![crpyto-implementaion-seq](./images/crpyto-implementaion-seq.png)
						
### Author

@yojkim
