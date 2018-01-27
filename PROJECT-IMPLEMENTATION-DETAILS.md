# Project Implementation Details



## Overview

It describes the important implementation decisions of the it-chain. Sample code for each detailed implementation can be found in the sample folder. 



## Table of Contents

1. [BlockChain](#BlockChain)
2. [SmartContract](#SmartContract)
3. [Communication](#Communication)



## BlockChain <a name="BlockChain"></a>



## SmartContract <a name="SmartContract"></a>



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

@Junbeomlee



## Crpyto

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
<br>

#### Signing process of data
![crpyto-implementaion-seq](./images/crpyto-implementaion-seq.png)

​									

### Author

@yojkim
