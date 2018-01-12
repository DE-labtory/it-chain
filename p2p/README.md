# P2P Network
## Features
### Heartbeat
Maintaining the list of nodes that participate in the network through checking if neighbors are alive.

#### Push `NodeTable`
Periodically, a node push its `NodeTable` to its neighbors. Neighbors are selected randomly from its `NodeTable`.

#### Pull `NodeTable`
##### Scenario 1: Joining the network.
When a new node tries to join the network with the address of the boot node(any node that already participates in the network is okay.), the new node will pull `NodeTable` from the boot node.

##### Scenario 2: Escaping the isolation.
If a node in the network doesn't receive any `NodeTable` from neighbors, the node tries to pull `NodeTable` from one of its neighbors.

### Consensus
If it-chain uses PBFT, it will need following additional functions in this `p2p` package.

#### Pull transactions from every node
When a node is ready to create a block, the node will pull all transactions that are not in a block using this function.

#### Push the vote to a leader node
After verifying the block individually in a single node, it will return the result to a leader node.

#### Push Block (newly created or verified) to all nodes
The leader node will use this function to propagate the newly-created or verified block to all nodes.

## Proposed structure for `P2P` package

```go
type NodeTable struct {
	// It should be created and maintained once. Maybe we can consider building a server type to hold this struct.
	nodeTable map[string]*Node
}

type Node struct {
	// attributes such as id, ipAddress, ...etc.
}

func (nt *NodeTable) Push (message interface{}, isGossipMode bool) {
	switch v := message.(type) {
		case SomeConcreteMessageType1:
			// For example, a message for Scenario#1
			// call some not-exported function to handle 
		case ...:
		...
	}
}

// Same as `Pull` function.
```

### Suggestion
* Consensus needs to be a separate package, which uses `Push` and `Pull` functions of `p2p` package to do features that are identified above.
* `p2p` package needs to focus on invoking `grpc` functions from nodes defined in `NodeTable` in terms of the separation of concerns.