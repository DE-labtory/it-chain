# Key Concept of Blockchain



## Overview

It describes the important concept of blockchains.

## Comparison of Blockchains

| Criteria    | Bitcoin                                  | Ethereum                     | HyperLedger                              |
| :---------- | :--------------------------------------- | :--------------------------- | :--------------------------------------- |
| Opensource? | Yes                                      | Yes                          | Yes                                      |
| Block time  | 10min                                    | 15seconds                    | Immediately                              |
| Consensus   | PoW                                      | PoW  -> PoS                  | Full-circle verification of the  correctness of a set of transactions comprising a block |
| Storage     | Transactions                             | Transactions + code          | Transactions                             |
| Description | Cryptocurrency based  on blockchain  platform | Generic  blockchain platform | Modular  blockchain platform             |
| Openness    | Private                                  | Public                       | Private  + Public                        |
| Transaction | Transfer of ownership                    | Message to send by account   | Chain Code Execution Message             |

#### Author

[Junbeomlee](https://github.com/junbeomlee)

## Blockchain

## SmartContract
Smart contracts are computer protocols that can facilitate, verify or enforce the negotiation or performance of a contract or make a contractual clause unnecessary. They usually have a User Interface and can emulate the logic of contractual clauses. They can execute the terms of a contract in an automated way. They can make contractual clauses partially or fully self-executing and self-enforcing.
Smart contracts are implemented using blockchain. Once a smart contract is created, it is placed in a blockchain.

### Advantages of Smart Contracts
- Smart Contracts eliminate the need of any intermediary like a broker, lawyer etc.
- The documents are encrypted in blockchain, which makes it much more secure. Also, the involved parties can be anonymous and maintain privacy.
- Usually a user has to spend lots of time for paperwork or to manually process documents. Smart contracts can automate the whole process, thereby saving time.
- As smart contracts eliminate the need of intermediaries, it saves costs involved in the whole process.
- As smart contracts are executed in an automated manner, it helps in avoiding errors that result from manual execution.

## Consensus

The use of a consensus algorithm is essential because anyone participating in a block chain can enter, change or delete data. The consensus algorithm can guarantee the reliability of the data stored in the block chain after a specific mechanism operation between authorized users. Thereby blockchain can be safely updated and maintained the state of the block chain, ensuring data integrity within the block chain.

|                               | PoW               | PoS     | PBFT                   |
| ----------------------------- | ----------------- | ------- | ---------------------- |
| **Major Blockchain Platform** | Bitcoin, Ethereum | Cardano | Hyperledger fabric 0.6 |

### PBFT

 PBFT is an algorithm introduced by Miguel Castro and Barbara Riskop in 1999. It is an algorithm designed for high-speed transaction processing and is capable of handling tens of millions of transactions per second. All participants in the network must be known in advance, and one of the participants will be the Primary. Send a request to all participants. The result for that request is aggregated and the block is committed using multiple values.

- Protocols for synchronizing state machines between multiple replicas
- Even though (N-1) / 3 participants try to fake at the same time, it can withstand.

![consensus-keyconecpt-pbft](./images/consensus-keyconecpt-pbft.png)

#### Process

1. A client sends a request to invoke a service operation to the primary (or leader)
2. The primary (or leader) multicasts the request to peers
3. Each peer executes the request and send a reply to the client
4. The client waits for f + 1 replies from different replicas with the same result; This is the result of the operation.

#### Overall behavior (4 peers)

![](https://github.com/bigpicturelabs/consensusPBFT/blob/master/pbft-consensus-behavior.jpg)

Definitions of each abbreviation in the diagram are;

- `m`: Request message object
- `c`: Client ID
- `t`: Timestamp
- `v`: View ID
- `n`: Sequence ID
- `i`: Peer(Node) ID
- `r`: Result of the request's operation

##### Why `count >= 2` ?

In the diagram, the peer change its state to `prepared` or `committed` when the `count` value, which is the number of verified messages from other peers, is larger than `2`.
Actually, the condition is `count >= 2*f` where `f` is the maximum number of faulty peers, which the network can tolerate. In this case, `f` is just `1`, so the condition is `count >= 2`. 

#####  What is the reply message?

Every node replies the result of the request's operation to the client individually. The client will collect these reply messages and if `f + 1` valid reply messages are arrived, the client will accept the result.

#### Sample implementation of PBFT

[bigpicturelabs/consensusPBFT](https://github.com/bigpicturelabs/consensusPBFT)

#### License

The copyright of overall behavior of pbft and the sample implementation of pbft are in bigpicturelab.

#### Author

[bitpicturelab](https://github.com/bigpicturelabs)

[Junbeomlee](https://github.com/junbeomlee)



### POW



### POS





