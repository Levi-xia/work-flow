package parser

import (
	"testing"
	"workflow/internal/core/store"
)

func TestParseProcessDefine(t *testing.T) {
	define, _ := store.GetProcessDefineStore().GetProcessDefine(100001)
	process, err := Parser2Process(define.Content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(process)
}
