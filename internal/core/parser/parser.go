package parser

import (
	"encoding/json"
	"time"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
	"workflow/internal/core/process"
)

type ProcessDefineParser struct {
	Code   string
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

type HttpActionParser struct {
	Url     string
	Method  string
	Timeout int
	Headers map[string]string
}

type ActionParser struct {
	ActionType constants.ActionType
	HttpAction HttpActionParser
	Params     map[string]interface{}
}

type Properties struct {
	Expr             string
	PreHook          *ActionParser
	PreInterceptors  []*ActionParser
	PostInterceptors []*ActionParser
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
			Code:  parser.Code,
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
				Code:  nodeParser.ID,
				Name: nodeParser.Text.Value,
			},
			Type:    nodeParser.Type,
			Inputs:  make([]*base.Edge, 0),
			Outputs: make([]*base.Edge, 0),
		}
		if nodeParser.Properties.PreHook != nil {
			node.PreHook = &base.Action{
				ActionType: nodeParser.Properties.PreHook.ActionType,
				HttpAction: base.HttpAction{
					Url:     nodeParser.Properties.PreHook.HttpAction.Url,
					Method:  nodeParser.Properties.PreHook.HttpAction.Method,
					Timeout: time.Duration(nodeParser.Properties.PreHook.HttpAction.Timeout) * time.Millisecond,
					Headers: nodeParser.Properties.PreHook.HttpAction.Headers,
				},
			}
		}
		for _, interceptor := range nodeParser.Properties.PreInterceptors {
			node.PreInterceptors = append(node.PreInterceptors, &base.Action{
				ActionType: interceptor.ActionType,
				HttpAction: base.HttpAction{
					Url:     interceptor.HttpAction.Url,
					Method:  interceptor.HttpAction.Method,
					Timeout: time.Duration(interceptor.HttpAction.Timeout) * time.Millisecond,
					Headers: interceptor.HttpAction.Headers,
				},
			})
		}
		for _, interceptor := range nodeParser.Properties.PostInterceptors {
			node.PostInterceptors = append(node.PostInterceptors, &base.Action{
				ActionType: interceptor.ActionType,
				HttpAction: base.HttpAction{
					Url:     interceptor.HttpAction.Url,
					Method:  interceptor.HttpAction.Method,
					Timeout: time.Duration(interceptor.HttpAction.Timeout) * time.Millisecond,
					Headers: interceptor.HttpAction.Headers,
				},
			})
		}
		nodeMap[nodeParser.ID] = node
		process.Nodes = append(process.Nodes, node)
	}

	// 构建所有边
	for _, edgeParser := range parser.Edges {
		edge := &base.Edge{
			Base: base.Base{
				Code:  edgeParser.ID,
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
