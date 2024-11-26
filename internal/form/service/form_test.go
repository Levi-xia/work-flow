package service

import (
	"testing"
	"workflow/internal/serctx"
)

func TestCreateFormDefine(t *testing.T) {
	serctx.InitServerContext()

	componentStructure := `[
  {
    "type": "input",
    "field": "days",
    "title": "请假天数",
    "info": "",
    "$required": "天数必填",
    "props": {
      "type": "number",
      "clearable": false
    },
    "_fc_id": "id_F6t3m3sfhg5zadc",
    "name": "ref_F103m3sfi4maagc",
    "display": true,
    "hidden": false,
    "_fc_drag_tag": "input"
  },
  {
    "type": "input",
    "field": "reason",
    "title": "请假理由",
    "info": "",
    "$required": "理由必填",
    "props": {
      "type": "textarea"
    },
    "_fc_id": "id_Fzcum3sfkexbakc",
    "name": "ref_F9ylm3sfkexbalc",
    "display": true,
    "hidden": false,
    "_fc_drag_tag": "textarea"
  }
]`

	formStructure := `{
  "form": {
    "inline": false,
    "hideRequiredAsterisk": false,
    "labelPosition": "right",
    "size": "default",
    "labelWidth": "125px"
  },
  "globalEvent": {
    "event_Feq4lui56zxbabc": {
      "label": "自定义",
      "deletable": false,
      "handle": "[[FORM-CREATE-PREFIX-function handle($inject){aaa;}-FORM-CREATE-SUFFIX]]"
    },
    "event_Feq4lui56zxbab2c": {
      "label": "自定义2",
      "handle": "[[FORM-CREATE-PREFIX-function(e){console.log(e)}-FORM-CREATE-SUFFIX]]"
    }
  },
  "globalData": {
    "data_Fk6dlui4k0xuabc": {
      "label": "自定义数据",
      "type": "static",
      "data": [
        1,
        2,
        3,
        4
      ]
    },
    "data_Fs1elui4kttlacc": {
      "action": "http://192.168.1.4:8081/",
      "deletable": false,
      "method": "GET",
      "headers": {},
      "data": {},
      "parse": "[[FORM-CREATE-PREFIX-function parse(res){return res.data;}-FORM-CREATE-SUFFIX]]",
      "onError": "[[FORM-CREATE-PREFIX-function onError(e){}-FORM-CREATE-SUFFIX]]",
      "label": "自定义接口数据",
      "type": "fetch"
    }
  },
  "resetBtn": {
    "show": false,
    "innerText": "重置"
  },
  "submitBtn": {
    "show": true,
    "innerText": "提交"
  },
  "globalVariable": {
    "var_Fppdlz6gytmzb1c": {
      "label": "token",
      "deletable": false,
      "handle": "[[FORM-CREATE-PREFIX-function(e,n){return e(\"$cookie.token\")||\"default Token\"}-FORM-CREATE-SUFFIX]]"
    }
  },
  "globalClass": {
    "cls_Fzmulzw3u0oib0c": {
      "label": "fff",
      "deletable": false,
      "style": {
        "color": "red"
      }
    }
  }
}`

	define, err := NewFormDefine("请假通用表单", "leave_form", formStructure, componentStructure)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(define.Meta.ToBo())
}

func TestCreateActionDefine2(t *testing.T) {
  serctx.InitServerContext()
	componentStructure := `[
  {
    "type": "input",
    "field": "remark",
    "title": "审批建议",
    "$required": false,
    "props": {
      "type": "textarea"
    },
    "_fc_id": "id_Fxzwm3x1djeoadc",
    "name": "ref_Fsbmm3x1djeoaec",
    "display": true,
    "hidden": false,
    "_fc_drag_tag": "textarea"
  }
]`
	formStructure := `
{
  "form": {
    "inline": false,
    "hideRequiredAsterisk": false,
    "labelPosition": "right",
    "size": "default",
    "labelWidth": "125px"
  },
  "globalEvent": {
    "event_Feq4lui56zxbabc": {
      "label": "自定义",
      "deletable": false,
      "handle": "[[FORM-CREATE-PREFIX-function handle($inject){aaa;}-FORM-CREATE-SUFFIX]]"
    },
    "event_Feq4lui56zxbab2c": {
      "label": "自定义2",
      "handle": "[[FORM-CREATE-PREFIX-function(e){console.log(e)}-FORM-CREATE-SUFFIX]]"
    }
  },
  "globalData": {
    "data_Fk6dlui4k0xuabc": {
      "label": "自定义数据",
      "type": "static",
      "data": [
        1,
        2,
        3,
        4
      ]
    },
    "data_Fs1elui4kttlacc": {
      "action": "http://192.168.1.4:8081/",
      "deletable": false,
      "method": "GET",
      "headers": {},
      "data": {},
      "parse": "[[FORM-CREATE-PREFIX-function parse(res){return res.data;}-FORM-CREATE-SUFFIX]]",
      "onError": "[[FORM-CREATE-PREFIX-function onError(e){}-FORM-CREATE-SUFFIX]]",
      "label": "自定义接口数据",
      "type": "fetch"
    }
  },
  "resetBtn": {
    "show": false,
    "innerText": "重置"
  },
  "submitBtn": {
    "show": true,
    "innerText": "提交"
  },
  "globalVariable": {
    "var_Fppdlz6gytmzb1c": {
      "label": "token",
      "deletable": false,
      "handle": "[[FORM-CREATE-PREFIX-function(e,n){return e(\"$cookie.token\")||\"default Token\"}-FORM-CREATE-SUFFIX]]"
    }
  },
  "globalClass": {
    "cls_Fzmulzw3u0oib0c": {
      "label": "fff",
      "deletable": false,
      "style": {
        "color": "red"
      }
    }
  }
}`
	define, err := NewFormDefine("审批通用表单", "approval_form", formStructure, componentStructure)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(define.Meta.ToBo())
}
