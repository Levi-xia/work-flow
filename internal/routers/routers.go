package routers

import (
	"workflow/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {

	// 拦截器
	interceptionRouter := r.Group("workflow/interception/v1")
	{
		interceptionRouter.GET("/sendSms", handler.SendSms)
	}

	// 动作路由
	actionRouter := r.Group("workflow/action/v1")
	{
		actionRouter.POST("createActionDefine", handler.CreateActionDefine)
	}

	// 表单路由
	formRouter := r.Group("workflow/form/v1")
	{
		formRouter.POST("/createFormDefine", handler.CreateFormDefine)
	}

	// 流程路由
	processRouter := r.Group("workflow/process/v1")
	{
		// 创建流程定义
		processRouter.POST("/createProcessDefine", handler.CreateProcessDefine)
		// 启动流程实例
		processRouter.POST("/startProcessInstance", handler.StartProcessInstance)
		// 执行任务
		processRouter.POST("/executeTask", handler.ExecuteTask)
	}
}
