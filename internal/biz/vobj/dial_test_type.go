// Package vobj is the value object package for the Sovereign service.
package vobj

type DialTestType string

const (
	DialTestTypePing DialTestType = "ping" // ping
	DialTestTypeCert DialTestType = "cert" // 证书
	DialTestTypePort DialTestType = "port" // 端口
	DialTestTypeHTTP DialTestType = "http" // http
)

func (t DialTestType) String() string {
	return string(t)
}

