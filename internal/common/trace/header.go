package trace

import (
	"context"
	"net/http"
	"net/url"
)

func NewHeader(c context.Context, trace *Trace) map[string]string {
	if trace == nil {
		trace = New(&http.Request{
			URL: &url.URL{
				Path: "/",
			},
		})
	}
	header := make(map[string]string)
	header[HeaderTraceID.K()] = trace.TraceId.V()
	header[HeaderSpanID.K()] = trace.SpanId.V()
	var parentSpanId string
	if trace.ParentSpanId != nil {
		parentSpanId = trace.ParentSpanId.V()
	} else {
		parentSpanId = RandSeq(6)
	}
	header[HeaderParentSpanID.K()] = parentSpanId
	return header
}

func getTraceV(h http.Header, k string) string {
	rsp := h.Get(k)
	if rsp == "" {
		rsp = RandSeq(6)
	}
	return rsp
}
