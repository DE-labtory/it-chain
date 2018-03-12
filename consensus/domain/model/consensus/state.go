package consensus

type State interface{
	Process()
}

type IdleState struct{

}

func (IdleState) Process(){

}

type PreprepareState struct{

}

func (PreprepareState) Process(){

}

type PrepareState struct{

}

func (PrepareState) Process(){

}


type CommitState struct{

}

type EndState struct{

}
