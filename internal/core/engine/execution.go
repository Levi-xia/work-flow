package engine

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"workflow/internal/core/base"
	"workflow/internal/core/constants"
	"workflow/internal/core/parser"
	"workflow/internal/core/process"
	"workflow/internal/core/service"
	"workflow/internal/core/store"
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
func (e *Execution) ExecuteTask(taskId int) error {
	var (
		task     *service.ProcessTask     = &service.ProcessTask{}
		instance *service.ProcessInstance = &service.ProcessInstance{}
		define   *service.ProcessDefine   = &service.ProcessDefine{}
		err      error
	)
	// 这里执行权限的判断 Todo

	// 获取store
	define.Store, instance.Store, task.Store = store.GetProcessDefineStore(), store.GetProcessInstanceStore(), store.GetProcessTaskStore()
	// 拿到定义
	if task.Meta, err = task.Store.GetProcessTask(taskId); err != nil {
		return err
	}
	// 获取instance
	if instance.Meta, err = instance.Store.GetProcessInstance(task.Meta.ProcessInstanceID); err != nil {
		return err
	}
	// 获取define
	if define.Meta, err = define.Store.GetProcessDefine(instance.Meta.ProcessDefineID); err != nil {
		return err
	}
	// 获取Process
	if e.Process, err = parser.Parser2Process(define.Meta.Content); err != nil {
		return err
	}
	instance.Define = define
	// task.Instance = instance
	e.Instance = instance

	// 任务、实例状态校验
	if instance.Meta.Status != constants.PROCESSINSTANCESTATUSDOING {
		return errors.New("instance status is not doing")
	}
	if task.Meta.Status != constants.PROCESSTASKSTATUSDOING {
		return errors.New("task status is not doing")
	}
	// 根据task信息获取node
	node, err := e.Process.GetNode(task.Meta.Code)
	if err != nil {
		return err
	}
	// instance公共变量赋值
	for k, v := range instance.Meta.Variables {
		if _, ok := e.Variables[k]; !ok {
			e.Variables[k] = v
		}
	}
	// 前置拦截器
	if node.PreInterceptors != nil {
		for _, interceptor := range node.PreInterceptors {
			if err := e.executeAction(interceptor); err != nil {
				return err
			}
		}
	}
	// 执行边的execute逻辑
	errCh := make(chan error)
	go func() {
		var wg sync.WaitGroup
		for _, edge := range node.Outputs {
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
	// 结束任务写库
	if err := task.Store.FinishProcessTask(task.Meta.ID); err != nil {
		return err
	}
	// 后置拦截器
	if node.PostInterceptors != nil {
		for _, interceptor := range node.PostInterceptors {
			if err := e.executeAction(interceptor); err != nil {
				return err
			}
		}
	}
	return nil
}

// 执行开始节点
func (e *Execution) executeStartNode(node *base.Node) error {
	if e.Define == nil {
		return errors.New("define is nil")
	}
	var err error
	// 创建实例
	if e.Instance, err = service.NewProcessInstance(e.Define, e.Variables, store.GetProcessInstanceStore()); err != nil {
		return err
	}
	return e.runOutputEdges(node.Outputs)
}

// 创建任务
func (e *Execution) createTask(node *base.Node) error {
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	if _, err := service.NewProcessTask(e.Instance,node.Code, node.Name, e.Variables); err != nil {
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
func (e *Execution) runOutputEdges(edges []*base.Edge) error {
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
	return e.runOutputEdges(node.Outputs)
}

// 结束实例
func (e *Execution) finishInstance() error {
	if e.Instance == nil {
		return errors.New("instance is nil")
	}
	return e.Instance.FinishProcessInstance()
}

// 执行节点动作
func (e *Execution) executeAction(action *base.Action) error {
	switch action.ActionType {
	case constants.HTTPCALLED:
		_, err := utils.HttpDo(action.HttpAction.Url, action.HttpAction.Method,
			utils.WithHeaders(action.HttpAction.Headers),
			utils.WithTimeout(action.HttpAction.Timeout))
		return err
	}
	return nil
}
