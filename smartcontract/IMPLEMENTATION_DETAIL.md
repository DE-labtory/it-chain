## SmartContract <a name="SmartContract"></a>

![smartContract-implementation-deploy](../images/smartContract-implementation-deploy.png)

SmartContract is stored on git repository and is executed by the smart contract service. After testing Smart Contract in a Docker-based virtual environment, it is reflected in the actual database.

- Git

  Each Smart Contract is stored as a Git Repository.

- Docker VM

  It is a virtual environment that executes smart contracts. After the smart contract and the world state db are copied to the Docker vm, they are executed and verified virtually.

- SmartContractService

  It is a service that manages git and Docker VM. After pushing and cloning the smart contract on the git, it copies the world state DB and smart contract to Docker VM and executes it.

  ​

#### Deploy Smart Contract Sequence Diagram

![smartContract-implementation-seq](../images/smartContract-implementation-seq.png)

The deployed user's repository is stored and managed in the Authenticated Smart Contract Repository as shown below.

| User <br />Repository <br />Path | Smart Contract <br />Repository <br />Path | Smart Contract File Path                 |
| -------------------------------- | ---------------------------------------- | ---------------------------------------- |
| A/a                              | {authenticated_git_id}/A_a               | It-chain/SmartContracts/A_a/{commit_hash} |
| B/b                              | {authenticated_git_id}/B_b               | It-chain/SmartContracts/B_b/{commit_hash} |
| C/c                              | {authenticated_git_id}/C_c               | It-chain/SmartContracts/C_c/{commit_hash} |

### Author
[@hackurity01](https://github.com/hackurity01)

[@hackurity01](https://github.com/hackurity01)
