package handler

import (
	"net/http"
	"workflow/internal/common"
	"workflow/internal/core/engine"
	"workflow/internal/core/parser"
	"workflow/internal/core/process"
	"workflow/internal/core/service"
	"workflow/internal/core/store"
	"workflow/internal/dto"

	"github.com/gin-gonic/gin"
)

// 创建流程定义
func CreateProcessDefine(c *gin.Context) {
	form := &dto.CreateProcessDefineRequest{}
	rsp := &common.Result{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ParamError, common.GetErrorMsg(form, err)))
		return
	}
	processDefine, err := service.NewProcessDefine(form.Content)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, rsp.Success(&dto.CreateProcessDefineResponse{
		ProcessDefineId: processDefine.Meta.ID,
		Code:            processDefine.Meta.Code,
		Name:            processDefine.Meta.Name,
		Version:         processDefine.Meta.Version,
	}))
}

// 启动流程实例
func StartProcessInstance(c *gin.Context) {
	form := &dto.StartProcessInstanceRequest{}
	rsp := &common.Result{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ParamError, common.GetErrorMsg(form, err)))
		return
	}
	defineMeta, err := store.GetProcessDefineStore().GetProcessDefine(form.ProcessDefineId)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	process, err := parser.Parser2Process(defineMeta.Content)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	define := &service.ProcessDefine{
		Meta:  defineMeta,
		Store: store.GetProcessDefineStore(),
	}
	if form.Variables == nil {
		form.Variables = make(map[string]interface{})
	}
	instance, err := service.NewProcessInstance(define, form.Variables)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	execution := &engine.Execution{
		Process:   process,
		Define:    define,
		Instance:  instance,
		Variables: form.Variables,
	}
	start, err := process.GetStart()
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	err = execution.ExecuteNode(start)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, rsp.Success(&dto.StartProcessInstanceResponse{
		ProcessInstanceId: instance.Meta.ID,
	}))
}

// 执行任务
func ExecuteTask(c *gin.Context) {
	form := &dto.ExecuteTaskRequest{}
	rsp := &common.Result{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ParamError, common.GetErrorMsg(form, err)))
		return
	}
	if form.Variables == nil {
		form.Variables = make(map[string]interface{})
	}
	var (
		task     *service.ProcessTask     = &service.ProcessTask{}
		instance *service.ProcessInstance = &service.ProcessInstance{}
		define   *service.ProcessDefine   = &service.ProcessDefine{}
		process  *process.Process
		err      error
	)
	// 获取store
	define.Store, instance.Store, task.Store = store.GetProcessDefineStore(), store.GetProcessInstanceStore(), store.GetProcessTaskStore()
	// 拿到定义
	if task.Meta, err = task.Store.GetProcessTask(form.TaskId); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	// 获取instance
	if instance.Meta, err = instance.Store.GetProcessInstance(task.Meta.ProcessInstanceID); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	// 获取define
	if define.Meta, err = define.Store.GetProcessDefine(instance.Meta.ProcessDefineID); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	// 获取Process
	if process, err = parser.Parser2Process(define.Meta.Content); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	instance.Define = define
	task.Instance = instance
	e := &engine.Execution{
		Process:   process,
		Define:    define,
		Instance:  instance,
		Variables: form.Variables,
	}
	err = e.ExecuteTask(task)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, rsp.Success(&dto.ExecuteTaskResponse{
		TaskId: task.Meta.ID,
	}))
}
