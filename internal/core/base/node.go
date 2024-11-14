package base

import (
	"workflow/internal/core/constants"
)

type Node struct {
	Base
	Type             constants.NodeType
	Inputs           []*Edge
	Outputs          []*Edge
	// 前置钩子（创建阻塞任务时执行）
	PreHook          *Action
	// 前置拦截器（实际执行节点前执行）
	PreInterceptors  []*Action
	// 后置拦截器（实际执行节点后执行）
	PostInterceptors []*Action
}
