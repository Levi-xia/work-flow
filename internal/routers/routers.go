package routers

import (
	"net/http"
	"workflow/internal/common"
	"workflow/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {

	r.NoRoute(func(c *gin.Context) {
		rsp := &common.Result{}
		c.JSON(http.StatusOK, rsp.Error(common.SystemError, "接口不存在"))
	})

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

	// 测试
	testRouter := r.Group("workflow/test/v1")
	{
		testRouter.GET("/token", func(ctx *gin.Context) {
			token, err := (&common.JwtService{}).CreateToken(common.AdminGuardName, "1000")
			ctx.JSON(200, gin.H{
				"token": token,
				"err":   err,
			})
		})
	}
}
