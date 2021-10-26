package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func StartSpanFromHttpRequest(r *http.Request, tracer opentracing.Tracer) opentracing.Span {
	spanCtx, _ := ExtractHttpSpan(r, tracer)
	return tracer.StartSpan("http-receive", ext.RPCServerOption(spanCtx))
}

func StartSpanFromGrpcRequest(ctx context.Context, tracer opentracing.Tracer) opentracing.Span {
	spanCtx, _ := ExtractGrpcSpan(ctx, tracer)
	return tracer.StartSpan("grpc-receive", ext.RPCServerOption(spanCtx))
}

func NewTracer(ctx context.Context, logger *zap.Logger, serviceName string) (opentracing.Tracer, error) {
	cfg := jconfig.Configuration{ServiceName: serviceName}
	tracer, closer, err := cfg.NewTracer(jconfig.Logger(NewJeagerLogger(logger)))
	if err != nil {
		return nil, fmt.Errorf("cannot init Jaeger tracer: %v", err)
	}

	go func() {
		<-ctx.Done()
		if err := closer.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()
	return tracer, nil
}
