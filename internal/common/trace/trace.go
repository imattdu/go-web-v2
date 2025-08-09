package trace

import (
	"net/http"
	"strings"
)

func New(req *http.Request) *Trace {
	return &Trace{
		TraceId: &entry{
			k: entryKTraceID,
			v: getTraceV(req.Header, HeaderTraceID.K()),
		},
		SpanId: &entry{
			k: entryKSpanID,
			v: getTraceV(req.Header, HeaderParentSpanID.K()),
		},
		ParentSpanId: &entry{
			k: entryKParentSpanID,
			v: "",
		},
		Uri: &entry{
			k: entryKURI,
			v: req.URL.RequestURI(),
		},
		From: &entry{
			k: entryKFrom,
			v: req.RemoteAddr,
		},
	}
}

func (t *Trace) Update(e entry, v string) {

}

func (t *Trace) UpdateParentSpanID() {
	t.ParentSpanId = &entry{
		k: entryKParentSpanID,
		v: RandSeq(6),
	}
}

func (t *Trace) Copy() *Trace {
	return &Trace{
		TraceId: &entry{
			k: entryKTraceID,
			v: t.TraceId.v,
		},
		SpanId: &entry{
			k: entryKSpanID,
			v: t.SpanId.v,
		},
		ParentSpanId: &entry{
			k: entryKParentSpanID,
			v: t.ParentSpanId.v,
		},
		Uri: &entry{
			k: entryKURI,
			v: t.Uri.v,
		},
		From: &entry{
			k: entryKFrom,
			v: t.From.v,
		},
	}
}

func (t *Trace) String() string {
	if t == nil {
		return ""
	}
	parts := make([]*entry, 0, 1)
	parts = append(parts, t.TraceId)
	parts = append(parts, t.SpanId)
	parts = append(parts, t.ParentSpanId)
	parts = append(parts, t.Uri)
	parts = append(parts, t.From)
	var sb strings.Builder
	for _, part := range parts {
		if part == nil {
			continue
		}
		sb.WriteString(part.K())
		sb.WriteString("=")
		sb.WriteString(part.v)
		sb.WriteString("||")
	}
	return sb.String()
}
