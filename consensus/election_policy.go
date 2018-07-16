package consensus

import "errors"

func Elect(parliament Parliament) ([]*Representative, error) {
	representative := make([]*Representative, 0)

	if !parliament.HasLeader() {
		return nil, errors.New("No Leader")
	}

	representative = append(representative, NewRepresentative(parliament.Leader.GetID()))

	for _, member := range parliament.Members {
		representative = append(representative, NewRepresentative(member.GetID()))
	}

	return representative, nil
}
