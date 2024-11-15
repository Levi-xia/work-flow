package main

import (
	"workflow/internal/cmd"
	"workflow/internal/common"
)

func main() {
	// 初始化字段验证器
	common.InitValidator()
	// 启动服务
	cmd.RunServer()
}
