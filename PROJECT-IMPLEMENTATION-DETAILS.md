# Project Implementation Details



## Overview

It describes the important implementation decisions of the it-chain. Sample code for each detailed implementation can be found in the sample folder. 



## Table of Contents

1. [BlockChain](#BlockChain)
2. [SmartContract](#SmartContract)
3. [Communication](#Communication)
4. [Crypto](#Crypto)
5. [Database](#DB)



## BlockChain <a name="BlockChain"></a>

![blockchain-implemeneation-logical](./images/blockchain-implemeneation-logical.png)

- BlockChain

  A blockchain is a continuously growing list of blocks, which are linked by a [hash](https://en.wikipedia.org/wiki/Cryptographic_hash_function) pointer as a link to a previous block.

- Block

  The block consists of a block header and block data, and the next block has a value obtained by hashing the block header for the block structure of the Ledger.
  The block header has the previous block hash value and the merge tree root hash value. The block data has a transaction list and has a merge tree To efficiently manage forgery and tampering of transactions.

- Transaction

  It is an atomic operation that performs Smart Contract. Transaction has TxData which contains ID of Peer (Node) to actually execute, transaction hash value that hashes transaction header, and contract contents.

- MerkleTree

  The Merkle Tree consists of a binary tree, and the leaf node is the hash value of the transactions in the transaction list of the block. The root node is a hash value representing the entire transaction that hashes the transaction hash value pair from the leaf node to the end. The merge tree can be immediately known at constant time via a hash root hash to see if the transaction's information has changed, since the next block has the value of the hash value of the block header and hashes the block header, the validity of all the transactions of the Ledger can be efficiently managed. And Merkle Path (a sibling node to the root node of tx), so that the validity of a particular transaction can be checked at log time.

  ![blockchain-implementation-merkletree](./images/blockchain-implementation-merkletree.png)

#### Author

[@emperorhan](https://github.com/emperorhan)

## SmartContract <a name="SmartContract"></a>

![smartContract-implementation-deploy](./images/smartContract-implementation-deploy.png)

SmartContract is stored on git repository and is executed by the smart contract service. After testing Smart Contract in a Docker-based virtual environment, it is reflected in the actual database.

- Git

  Each Smart Contract is stored as a Git Repository.

- Docker VM

  It is a virtual environment that executes smart contracts. After the smart contract and the world state db are copied to the Docker vm, they are executed and verified virtually.

- SmartContractService

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

Crypto signs and verifies the data used in the block-chain platform and manages the keys used in the process. *it-chain* supports `RSA` and `ECDSA` encryption method.

- KeyGenerator

  The node generates a key that matches the encryption scheme that you want to use for signing.

- KeyManager

  Stores the generated key, and loads the stored key.

- Signer

  Performs data signature.

- Verifier

  Verify the signed data.

- KeyUtils

  Perform the necessary processing tasks in the process of storing and loading the key such as converting key data to `PEM` file.

- Key

  Provides attribute values related to the interface of key data required in the signature or verification process.

<br>

### Signing process of data
![crpyto-implementaion-seq](./images/crpyto-implementaion-seq.png)
​						
### Author

[@yojkim](https://github.com/yojkim)

## Database <a name="DB"></a>
Blockchain can be stored in multiple types of database depend on configuration. For now leveldbhepler and fileheper functions are added. Basic DB implementation uses leveldb. Blocks can be retrieved by block hash and block number and transaction ID. Transactions can be retrieved by transaction ID.
If you want to use other database, implement it under blockchainleveldb and edit blockchain_db_interface.

### Related config
Database config is defined in config.yaml as database section

- type

  Type of database. For now only leveldb is supported and little helper function for file is supported.

- leveldb

  Configuration for leveldb.

  | key          | description                              |
  | ------------ | ---------------------------------------- |
  | default_path | If no other path for leveldb is provided, leveldb data is stored in this path |

### LevelDB
Blocks are totally stored in key-value storage leveldb.

- Blocks

  Blocks are serialized to JSON and saved in leveldb. For key block hash and block number are used.
  Last block and unconfirmed block are saved for recover.

- Transactions

  Also transactions are serialized and saved in leveldb. Basically all transactions are saved together block.
  For indexing, block hash that transaction belongs to also saved. Transaction ID is used for key.

| DB name      | Key            | Value                  | Description                              |
| ------------ | -------------- | ---------------------- | ---------------------------------------- |
| block_hash   | BlockHash      | Serialized Block       | Save block using blockhash               |
| block_number | BlockNumber    | Serialized Block       | Save block using block's number          |
| transaction  | Transaction ID | Serialized Transaction | Save transactions                        |
| util         | Predefined Key | Depends on Key         | Save last block, unconfirmed block, blockhash of block that transaction is stored, ... |

### File
Block's metadata is saved in leveldb or other key-value database. Block body is saved in file.

- Blocks

  Block's metadata is serialized to JSON and saved in leveldb. Block body data is written into file.

- Transactions

  Transaction data is stored in file. For finding, information of the file is stored in key-value database using transaction ID as key.

### Author

[@luke9407](https://github.com/luke9407)