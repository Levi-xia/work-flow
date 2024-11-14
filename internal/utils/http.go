package utils

import (
	"net/http"
	"time"
)

type HttpOption func(*httpRequest)

type httpRequest struct {
	client  *http.Client
	url     string
	method  string
	headers map[string]string
	timeout time.Duration
}

// 设置超时时间选项
func WithTimeout(timeout time.Duration) HttpOption {
	return func(r *httpRequest) {
		r.timeout = timeout
	}
}

// 设置请求头选项
func WithHeaders(headers map[string]string) HttpOption {
	return func(r *httpRequest) {
		r.headers = headers
	}
}

// 发送HTTP请求
func HttpDo(url string, method string, opts ...HttpOption) (*http.Response, error) {
	// 默认配置
	req := &httpRequest{
		client:  &http.Client{},
		url:     url,
		method:  method,
		headers: make(map[string]string),
		timeout: 30 * time.Second,
	}

	// 应用选项
	for _, opt := range opts {
		opt(req)
	}

	// 设置超时
	req.client.Timeout = req.timeout

	// 创建请求
	httpReq, err := http.NewRequest(req.method, req.url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range req.headers {
		httpReq.Header.Set(k, v)
	}

	// 发送请求
	return req.client.Do(httpReq)
}
