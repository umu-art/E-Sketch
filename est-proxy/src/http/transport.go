package http

import (
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"net/http"
)

type TransportWithTraceparentHeaders struct {
	childName string
}

func NewTransportWithTraceparentHeaders(childName string) *TransportWithTraceparentHeaders {
	return &TransportWithTraceparentHeaders{childName}
}

func (t *TransportWithTraceparentHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	traceContext := apm.TransactionFromContext(ctx).TraceContext()

	span, _ := apm.StartSpan(ctx, t.childName+": "+req.URL.Path, "service")
	defer span.End()

	// Set traceparent
	traceparent := apmhttp.FormatTraceparentHeader(traceContext)

	req.Header.Add("trace_id", traceparent)
	req.Header.Add("parent_span_id", traceparent)

	return http.DefaultTransport.RoundTrip(req)

}
