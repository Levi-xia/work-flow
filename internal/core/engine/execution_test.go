package engine

import (
	"testing"
	"workflow/internal/core/parser"
	"workflow/internal/core/service"
	"workflow/internal/core/store"
)

// 测试完整流程
func TestProcess(t *testing.T) {
	defineMeta, _ := store.GetProcessDefineStore().GetProcessDefine(100001)
	process, err := parser.Parser2Process(defineMeta.Content)
	if err != nil {
		t.Fatal(err)
	}
	execution := &Execution{
		Process: process,
		Define: &service.ProcessDefine{
			Meta:  defineMeta,
			Store: store.GetProcessDefineStore(),
		},
		Variables: map[string]interface{}{
			"days": 5,
		},
	}
	start, err := process.GetStart()
	if err != nil {
		t.Error(err)
	}
	err = execution.ExecuteNode(start)
	if err != nil {
		t.Error(err)
	}
}

func TestExecuteTask(t *testing.T) {
	execution := &Execution{
		Variables: map[string]interface{}{
			"days": 5,
		},
	}
	err := execution.ExecuteTask(100007)
	if err != nil {
		t.Error(err)
	}
}
