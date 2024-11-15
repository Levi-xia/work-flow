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
	define := &service.ProcessDefine{
		Meta:  defineMeta,
		Store: store.GetProcessDefineStore(),
	}
	variables := map[string]interface{}{
		"days": 5,
	}
	instance, err := service.NewProcessInstance(define, variables, store.GetProcessInstanceStore())
	if err != nil {
		t.Fatal(err)
	}
	execution := &Execution{
		Process:   process,
		Define:    define,
		Instance:  instance,
		Variables: variables,
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
}
