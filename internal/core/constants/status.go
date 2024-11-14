package constants

type ProcessInstanceStatus string

const (
	PROCESSINSTANCESTATUSDOING  ProcessInstanceStatus = "doing"
	PROCESSINSTANCESTATUSFINISH ProcessInstanceStatus = "finish"
	PROCESSINSTANCESTATUSCANCEL ProcessInstanceStatus = "cancel"
)


type ProcessTaskStatus string

const (
	PROCESSTASKSTATUSDOING  ProcessTaskStatus = "doing"
	PROCESSTASKSTATUSFINISH ProcessTaskStatus = "finish"
	PROCESSTASKSTATUSFAIL   ProcessTaskStatus = "fail"
)

