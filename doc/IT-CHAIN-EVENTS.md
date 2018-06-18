<!-- toc -->
# IT CHAIN EVENTS
Types of events(which is kind of AMQP messages) related to certain aggregate root and their usages

## Connection
| Event                       | Publisher    | Consumer     | Desc                                                    |
| --------------------------- | ------------ | ------------ | ------------------------------------------------------- |
| ConnectionCreatedEvent      | Gateway, P2P | Gateway, P2P | Events which is occured when connection is connected    |
| ConnectionDisconnectedEvent | Gateway, P2P | Gateway, P2P | Events which is occured when connection is disconnected |

## Leader
| Event              | Publisher | Consumer | Desc                                           |
| ------------------ | --------- | -------- | ---------------------------------------------- |
| LeaderUpdatedEvent | P2P       | P2P      | Events which is occured when leader is updated |

## Node
| Event            | Publisher | Consumer | Desc                                         |
| ---------------- | --------- | -------- | -------------------------------------------- |
| NodeCreatedEvent | P2P       | P2P      | Events which is occured when node is created |



---
### AUTHOR
Namhoon Lee(@frontalnh)