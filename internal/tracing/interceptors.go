package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

func InjectGrpcSpan(ctx context.Context, span opentracing.Span) (context.Context, error) {
	meta := make(map[string]string, 0)
	err := span.Tracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(meta))
	if err != nil {
		return nil, err
	}

	md := make([]string, 0)
	for k, v := range meta {
		md = append(md, k)
		md = append(md, v)
	}

	ct := metadata.AppendToOutgoingContext(ctx, md...)
	return ct, nil
}

func ExtractHttpSpan(r *http.Request, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

func ExtractGrpcSpan(ctx context.Context, tracer opentracing.Tracer) (opentracing.SpanContext, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("error duting grpc span extracting")
	}

	data := make(map[string]string, 0)
	for k, v := range meta {
		data[k] = v[0]
	}

	return tracer.Extract(
		opentracing.TextMap,
		opentracing.TextMapCarrier(data),
	)
}
