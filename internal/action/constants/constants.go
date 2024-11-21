package constants

import "net/http"

type ParamsStructType string

const (
	ParamsStructTypeString ParamsStructType = "string"
	ParamsStructTypeInt    ParamsStructType = "int"
	ParamsStructTypeBool   ParamsStructType = "bool"
	ParamsStructTypeFloat  ParamsStructType = "float"
	ParamsStructTypeObject ParamsStructType = "object"
)

type HttpMethod string

const (
	HttpMethodGet  HttpMethod = http.MethodGet
	HttpMethodPost HttpMethod = http.MethodPost
)
