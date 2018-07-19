package adapter

type Publish func(exchange string, topic string, data interface{}) (err error)

type PropagateService struct {
	publish Publish
}

func NewPropagateService(publish Publish) *PropagateService {
	return &PropagateService{
		publish: publish,
	}
}
