package main

import (
	"workflow/config"
	"workflow/internal/cmd"
	"workflow/internal/common"
	"workflow/internal/env"
	"workflow/internal/serctx"
)

func main() {
	// 初始化环境
	env.InitEnv()
	// 初始化配置
	config.InitConfig()
	// 初始化字段验证器
	common.InitValidator()
	// 初始化服务
	serctx.InitServerContext()
	// 启动服务
	cmd.RunServer()
}
