package trace

const (
	entryKTraceID      = "trace_id"
	entryKSpanID       = "span_id"
	entryKParentSpanID = "parent_span_id"
	entryKURI          = "uri"

	entryKFrom = "from"
)

var (
	HeaderTraceID = &entry{
		k: "Trace-Id",
	}
	HeaderSpanID = &entry{
		k: "Span-Id",
	}
	HeaderParentSpanID = &entry{
		k: "Parent-Span-Id",
	}
)
