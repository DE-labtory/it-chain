package api_test

//todo make node api test
//todo make fake dependencies 1. eventRepository 2. messageDispatcher 3. nodeRepository
//todo make test map
//todo test continue

//func TestNodeApi_UpdateNodeList(t *testing.T) {
//
//	tests := map[string]struct {
//		input []p2p.Node
//		err   error
//	}{
//		"empty node list test": {
//			input: []p2p.Node{},
//			err:   ErrEmptyNodeList,
//		},
//	}
//
//	publisher := messaging.Publisher(func(exchange string, topic string, data interface{}) error {
//		return nil
//	})
//
//	myInfo := p2p.NewNode(conf.GetConfiguration().Common.NodeIp, p2p.NodeId{Id: "1"})
//
//	nodeApi := SetupNodeApi("test", "test", "test", publisher, myInfo)
//
//	for testName, test := range tests {
//		t.Logf("running test case %s", testName)
//		err := nodeApi.UpdateNodeList(test.input)
//		assert.Equal(t, err, test.err)
//	}
//}
//
//func SetupNodeApi(nodeRepoPath string, eventRepoPath string, url string, publisher messaging.Publisher, myInfo *p2p.Node) *NodeApi {
//
//	event := p2p.LeaderUpdatedEvent{}
//	store := leveldb2.NewEventStore(eventRepoPath, leveldb2.NewSerializer(event))
//	client := rabbitmq.Connect(url)
//
//	nodeRepository := leveldb.NewNodeRepository(nodeRepoPath)
//	eventRepository := midgard.NewRepo(store, client)
//	messageDispatcher := messaging.NewMessageDispatcher(publisher)
//
//	nodeApi := NewNodeApi(nodeRepository, *eventRepository, messageDispatcher)
//
//	return nodeApi
//}
