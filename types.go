package zlog

import (
	"fmt"
	"net/http"
)

type ResponseLog struct {
	Request    *http.Request
	StatusCode int    `json:"staus_code"`
	Body       string `json:"body"`
	Header     string `json:"header"`
}

func (rl ResponseLog) DumpResponse() string {
	res := ""
	res += fmt.Sprintf("HTTP/%d.%d %d %s\r\n", rl.Request.ProtoMajor, rl.Request.ProtoMinor,
		rl.StatusCode, http.StatusText(rl.StatusCode))
	res += rl.Header
	res += "\r\n"
	res += rl.Body
	return res
}

// headerOperation represents an operation on the header
type headerOperation func(http.Header)

type ResponseProxyWriter struct {
	writer       http.ResponseWriter
	Body         []byte
	Code         int
	ops          []headerOperation
	SourceHeader http.Header
	wroteHeader  bool
}

func (this *ResponseProxyWriter) Header() http.Header {
	return this.writer.Header()
}

func (this *ResponseProxyWriter) Write(bytes []byte) (int, error) {
	if !this.wroteHeader {
		this.WriteHeader(http.StatusOK)
	}
	this.Body = append(this.Body, bytes[0:len(bytes)]...)
	return this.writer.Write(bytes)
}

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

func (this *ResponseProxyWriter) WriteHeader(i int) {
	if this.wroteHeader {
		return
	}
	this.wroteHeader = true

	h := this.writer.Header()
	this.SourceHeader = cloneHeader(h)

	// perform our revisions
	for _, op := range this.ops {
		op(h)
	}

	this.Code = i
	this.writer.WriteHeader(i)
}

func (this *ResponseProxyWriter) delHeader(key string) {
	// remove the existing one if any
	this.writer.Header().Del(key)

	// register a future deletion
	this.ops = append(this.ops, func(h http.Header) {
		h.Del(key)
	})
}

func NewRespProxyWriter(w http.ResponseWriter) *ResponseProxyWriter {
	return &ResponseProxyWriter{
		writer: w,
		Body:   []byte{},
	}
}
