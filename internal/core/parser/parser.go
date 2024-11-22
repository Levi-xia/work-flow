package parser

import (
	"encoding/json"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
	"workflow/internal/core/process"
)

type ProcessDefineParser struct {
	Code  string
	Name  string
	Nodes []NodeParser
	Edges []EdgeParser
}

type NodeParser struct {
	ID         string
	Type       constants.NodeType
	Properties Properties
	Text       Text
}

type EdgeParser struct {
	ID           string
	SourceNodeId string
	TargetNodeId string
	Properties   Properties
	Text         Text
}

type Properties struct {
	Expr             string `json:"expr"`
	FormID           int    `json:"form_id"`
	PreHooks         []int  `json:"pre_hooks"`
	PreInterceptors  []int  `json:"pre_interceptors"`
	PostInterceptors []int  `json:"post_interceptors"`
}

type Text struct {
	Value string
}

func Parser2Process(content string) (*process.Process, error) {
	var parser ProcessDefineParser
	err := json.Unmarshal([]byte(content), &parser)
	if err != nil {
		return nil, err
	}
	// 循环所有的点和线，找到对应的source和target关系，构建process
	process := &process.Process{
		Base: base.Base{
			Code: parser.Code,
			Name: parser.Name,
		},
		Nodes: make([]*base.Node, 0),
		Edges: make([]*base.Edge, 0),
	}
	// 先构建所有节点
	nodeMap := make(map[string]*base.Node)
	for _, nodeParser := range parser.Nodes {
		node := &base.Node{
			Base: base.Base{
				Code: nodeParser.ID,
				Name: nodeParser.Text.Value,
			},
			FormID:  nodeParser.Properties.FormID,
			Type:    nodeParser.Type,
			Inputs:  make([]*base.Edge, 0),
			Outputs: make([]*base.Edge, 0),
		}
		if len(nodeParser.Properties.PreHooks) > 0 {
			node.PreHooks = nodeParser.Properties.PreHooks
		}
		if len(nodeParser.Properties.PreInterceptors) > 0 {
			node.PreInterceptors = nodeParser.Properties.PreInterceptors
		}
		if len(nodeParser.Properties.PostInterceptors) > 0 {
			node.PostInterceptors = nodeParser.Properties.PostInterceptors
		}
		nodeMap[nodeParser.ID] = node
		process.Nodes = append(process.Nodes, node)
	}

	// 构建所有边
	for _, edgeParser := range parser.Edges {
		edge := &base.Edge{
			Base: base.Base{
				Code: edgeParser.ID,
				Name: edgeParser.Text.Value,
			},
			Expr: edgeParser.Properties.Expr,
		}

		// 找到source和target节点
		sourceNode := nodeMap[edgeParser.SourceNodeId]
		targetNode := nodeMap[edgeParser.TargetNodeId]

		edge.Source = sourceNode
		edge.Target = targetNode

		// 建立双向关联
		sourceNode.Outputs = append(sourceNode.Outputs, edge)
		targetNode.Inputs = append(targetNode.Inputs, edge)

		process.Edges = append(process.Edges, edge)
	}
	return process, nil
}
