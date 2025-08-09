package trace

type entry struct {
	k string
	v string
}

type Trace struct {
	TraceId      *entry `json:"trace_id"`
	SpanId       *entry `json:"span_id"`
	ParentSpanId *entry `json:"parent_span_id"`
	Uri          *entry `json:"uri"`
	From         *entry `json:"from"`
}
