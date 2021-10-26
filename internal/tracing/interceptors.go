package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
)

func InjectHttpSpan(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func InjectGrpcSpan(span opentracing.Span, ctx context.Context) (context.Context, error) {
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

func ExtractHttpSpan(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

func ExtractGrpcSpan(tracer opentracing.Tracer, ctx context.Context) (opentracing.SpanContext, error) {
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
