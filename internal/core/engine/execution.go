package engine

import (
	"errors"
	"fmt"
	"github.com/caibirdme/yql"
	"log"
	"sync"
	actionService "workflow/internal/action/service"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
	"workflow/internal/core/process"
	"workflow/internal/core/service"
	formBo "workflow/internal/form/bo"
	formService "workflow/internal/form/service"
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
	// 把instance的变量赋值到Variables，但是Variables已经存在的变量不覆盖
	for k, v := range e.Instance.Meta.Variables {
		if _, exists := e.Variables[k]; !exists {
			e.Variables[k] = v
		}
	}
	// 前置拦截器
	if err := actionService.ExecuteActions(node.PreInterceptors, task.Meta.ID, e.Variables); err != nil {
		return err
	}
	// 结束任务写库
	if err := service.FinishProcessTask(task.Meta.ID, e.Variables); err != nil {
		return err
	}
	// 更新instance的变量
	for k, v := range e.Variables {
		e.Instance.Meta.Variables[k] = v
	}
	if err := service.UpdateProcessInstanceVariables(e.Instance.Meta.ID, e.Instance.Meta.Variables); err != nil {
		return err
	}
	// 更新表单状态
	if task.Meta != nil && task.Meta.FormInstanceID > 0 {
		if err := formService.UpdateFormInstanceFormData(task.Meta.FormInstanceID, e.Variables); err != nil {
			return err
		}
	}
	// 后置拦截器
	if err := actionService.ExecuteActions(node.PostInterceptors, task.Meta.ID, e.Variables); err != nil {
		return err
	}
	// 执行边的execute逻辑
	if err := e.executeEdges(node.Outputs); err != nil {
		return err
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
	var (
		processTask *service.ProcessTask
		err         error
	)
	// 创建表单
	var formInstance *formService.FormInstance
	if node.FormID > 0 {
		if formInstance, err = formService.NewFormInstance(node.FormID, []*formBo.FormData{}); err != nil {
			return err
		}
	}
	if processTask, err = service.NewProcessTask(e.Instance, node.Code, node.Name, formInstance.Meta.ID, map[string]interface{}{}); err != nil {
		return err
	}
	// 执行前置拦截器
	if err := actionService.ExecuteActions(node.PreHooks, processTask.Meta.ID, e.Variables); err != nil {
		return err
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
				// 这里应该有问题，e被一个分支更新可能影响其他分支流程，可以考虑复制一份e执行边
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
