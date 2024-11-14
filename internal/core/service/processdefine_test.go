package service

import (
	"testing"
	"workflow/internal/core/store"
)

func TestCreateProcessDefine(t *testing.T) {
	str := `{
		"code": "leave",
		"name": "请假",
		"nodes": [
			{
				"id": "start",
				"type": "start",
				"properties": {
	
				},
				"text": {
					"value": "开始"
				}
			},
			{
				"id": "apply",
				"type": "task",
				"properties": {
	
				},
				"text": {
					"value": "请假申请"
				}
			},
			{
				"id": "approveDept",
				"type": "task",
				"properties": {
	
				},
				"text": {
					"value": "部门领导审批"
				}
			},
			{
				"id": "approveBoss",
				"type": "task",
				"properties": {
	
				},
				"text": {
					"value": "公司领导审批"
				}
			},
			{
				"id": "end",
				"type": "end",
				"properties": {
	
				},
				"text": {
					"value": "结束"
				}
			}
		],
		"edges": [
			{
				"id": "edge_1",
				"type": "transition",
				"sourceNodeId": "start",
				"targetNodeId": "apply",
				"properties": {
	
				}
			},
			{
				"id": "edge_2",
				"type": "transition",
				"sourceNodeId": "apply",
				"targetNodeId": "approveDept",
				"properties": {
	
				}
			},
			{
				"id": "edge_3",
				"type": "transition",
				"sourceNodeId": "approveDept",
				"targetNodeId": "approveBoss",
				"properties": {
					"expr": "days > 3"
				},
				"text": {
					"value": "请假天数大于3"
				}
			},
			{
				"id": "edge_4",
				"type": "transition",
				"sourceNodeId": "approveBoss",
				"targetNodeId": "end",
				"properties": {
	
				}
			},
			{
				"id": "edge_5",
				"type": "transition",
				"sourceNodeId": "approveDept",
				"targetNodeId": "end",
				"properties": {
					"expr": "days <= 3"
				},
				"text": {
					"value": "请假天数小于等于3"
				}
			}
		]
	}`
	processDefine, err := NewProcessDefine("leave", "请假", str, store.GetProcessDefineStore())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(processDefine)
}
