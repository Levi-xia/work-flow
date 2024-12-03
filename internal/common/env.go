package common

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"path/filepath"
)

// 全局运行环境
var ENV string

const (
	ENV_DEV  = "dev"
	ENV_PROD = "prod"
)

// 初始化运行环境
func InitEnv() {
	workPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("get work path error: %v\n", err)
	}
	serviceClusterFile := filepath.Join(workPath, "./.deploy/service.cluster.txt")

	f, err := os.Open(serviceClusterFile)
	if err != nil {
		log.Fatalf("open service.cluster.txt error: %v\n", err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if bytes.Equal(line, []byte{}) {
			continue
		} else if err != nil {
			log.Fatalf("read service.cluster.txt error: %v\n", err)
		}
		ENV = string(bytes.TrimSpace(line))
		break
	}
	fmt.Printf("Init global environment is：%s\n", ENV)

	// 设置gin运行模式
	if ENV == ENV_DEV {
		gin.SetMode(gin.DebugMode)
	} else if ENV == ENV_PROD {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Fatalf("ENV is not valid, ENV: %s\n", ENV)
	}
}
