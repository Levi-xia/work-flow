package process

import (
	"errors"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
)

type Process struct {
	base.Base
	Nodes []*base.Node
	Edges []*base.Edge
}

func (this *Process) GetStart() (*base.Node, error) {
	for _, node := range this.Nodes {
		if node.Type == constants.STARTNODETYPE {
			return node, nil
		}
	}
	return nil, errors.New("no start node")
}

func (this *Process) GetNode(code string) (*base.Node, error) {
	for _, node := range this.Nodes {
		if node.Code == code {
			return node, nil
		}
	}
	return nil, errors.New("node not found")
}
