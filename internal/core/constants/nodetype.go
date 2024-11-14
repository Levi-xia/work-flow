package constants

type NodeType string

const (
	STARTNODETYPE NodeType = "start"
	ENDNODETYPE   NodeType = "end"
	TASKNODETYPE  NodeType = "task"
	JOINNODETYPE  NodeType = "join"
)