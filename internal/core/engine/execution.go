package engine

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
	"workflow/internal/core/process"
	"workflow/internal/core/service"
	"workflow/internal/utils"

	"github.com/caibirdme/yql"
)

type Execution struct {
	// 流程对象
	Process *process.Process
	// 定义
	Define *service.ProcessDefine
	// 实例
	Instance *service.ProcessInstance
	// 变量
	Variables map[string]interface{}
	// 执行人ID
	UserID int
}

// 执行节点
func (e *Execution) ExecuteNode(node *base.Node) error {
	switch node.Type {
	case constants.STARTNODETYPE:
		log.Printf("start node %s execute", node.Code)
		return e.executeStartNode(node)
	case constants.ENDNODETYPE:
		log.Printf("end node %s execute", node.Code)
		return e.finishInstance()
	case constants.TASKNODETYPE:
		log.Printf("task node %s execute", node.Code)
		return e.createTask(node)
	case constants.JOINNODETYPE:
		log.Printf("join node %s execute", node.Code)
		return e.executeJoinNode(node)
	}
	return nil
}

// 执行任务
func (e *Execution) ExecuteTask(task *service.ProcessTask) error {
	if e.Define == nil {
		return errors.New("define is nil")
	}
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	// 这里执行权限的判断 Todo

	// 任务、实例状态校验
	if e.Instance.Meta.Status != constants.PROCESSINSTANCESTATUSDOING {
		return errors.New("instance status is not doing")
	}
	// 任务状态校验
	if task.Meta.Status != constants.PROCESSTASKSTATUSDOING {
		return errors.New("task status is not doing")
	}
	// 根据task信息获取node
	node, err := e.Process.GetNode(task.Meta.Code)
	if err != nil {
		return err
	}
	// task的变量赋值
	for k, v := range task.Meta.Variables {
		if _, ok := e.Variables[k]; !ok {
			e.Variables[k] = v
		}
	}
	// 前置拦截器
	if node.PreInterceptors != nil {
		if err := e.executeActions(node.PreInterceptors); err != nil {
			return err
		}
	}
	// 结束任务写库
	if err := service.FinishProcessTask(task.Meta.ID, e.Variables); err != nil {
		return err
	}
	// 执行边的execute逻辑
	if err := e.executeEdges(node.Outputs); err != nil {
		return err
	}
	// 更新instance的变量
	for k, v := range e.Variables {
		e.Instance.Meta.Variables[k] = v
	}
	if err := service.UpdateProcessInstanceVariables(e.Instance.Meta.ID, e.Instance.Meta.Variables); err != nil {
		return err
	}
	// 后置拦截器
	if node.PostInterceptors != nil {
		if err := e.executeActions(node.PostInterceptors); err != nil {
			return err
		}
	}
	return nil
}

// 执行开始节点
func (e *Execution) executeStartNode(node *base.Node) error {
	return e.executeEdges(node.Outputs)
}

// 创建任务
func (e *Execution) createTask(node *base.Node) error {
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	if e.Variables == nil {
		e.Variables = make(map[string]interface{})
	}
	// 把instance的变量赋值到Variables，但是Variables已经存在的变量不覆盖
	for k, v := range e.Instance.Meta.Variables {
		if _, exists := e.Variables[k]; !exists {
			e.Variables[k] = v
		}
	}
	if _, err := service.NewProcessTask(e.Instance, node.Code, node.Name, e.Variables); err != nil {
		return err
	}
	// 执行前置拦截器
	if node.PreHook != nil {
		if err := e.executeAction(node.PreHook); err != nil {
			return err
		}
	}
	return nil
}

// 批量执行输出边
func (e *Execution) executeEdges(edges []*base.Edge) error {
	errCh := make(chan error)
	go func() {
		var wg sync.WaitGroup
		for _, edge := range edges {
			wg.Add(1)
			go func(edge *base.Edge) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errCh <- fmt.Errorf("panic in edge execution: %v", r)
					}
				}()
				if err := e.executeEdge(edge); err != nil {
					errCh <- err
				}
			}(edge)
		}
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// 执行边
func (e *Execution) executeEdge(edge *base.Edge) error {
	if edge.Expr == "" {
		return e.ExecuteNode(edge.Target)
	}
	match, err := yql.Match(edge.Expr, e.Variables)
	if err != nil {
		return err
	}
	if !match {
		return nil
	}
	return e.ExecuteNode(edge.Target)
}

// 执行合并节点
func (e *Execution) executeJoinNode(node *base.Node) error {
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	runningTasks, err := service.GetRunningTasks(e.Instance)
	if err != nil {
		return err
	}
	if len(runningTasks) >= 0 {
		return nil
	}
	return e.executeEdges(node.Outputs)
}

// 结束实例
func (e *Execution) finishInstance() error {
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	return service.FinishProcessInstance(e.Instance.Meta.ID)
}

// 执行节点动作
func (e *Execution) executeAction(action *base.Action) error {
	if action.Params == nil {
		action.Params = make(map[string]interface{})
	}
	// 将Variables中的变量复制到Params,但不覆盖已有值
	for k, v := range e.Variables {
		if _, exists := action.Params[k]; !exists {
			action.Params[k] = v
		}
	}
	switch action.ActionType {
	case constants.HTTPCALLED:
		/*这里暂时不处理返回内容，如果要处理，对e进行修改，需要加锁
		同时e的变量被修改，可能会导致多个action执行时，变量内容不同，
		使用action能力需注意*/
		_, err := utils.HttpDo(action.HttpAction.Url, action.HttpAction.Method, action.Params,
			utils.WithHeaders(action.HttpAction.Headers),
			utils.WithTimeout(action.HttpAction.Timeout))
		return err
	}
	return nil
}

// 批量执行action
func (e *Execution) executeActions(actions []*base.Action) error {
	errCh := make(chan error)
	go func() {
		var wg sync.WaitGroup
		for _, action := range actions {
			wg.Add(1)
			go func(action *base.Action) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errCh <- fmt.Errorf("panic in action execution: %v", r)
					}
				}()
				if err := e.executeAction(action); err != nil {
					errCh <- err
				}
			}(action)
		}
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}
