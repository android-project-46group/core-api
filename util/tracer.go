package util

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	agentHost = "127.0.0.1:5775"
)

func NewJaegerTracer(service string) (opentracing.Tracer, io.Closer, error) {
	//nolint:exhaustivestruct,exhaustruct
	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: time.Second,
			LocalAgentHostPort:  agentHost,
		},
	}

	//nolint:wrapcheck
	return cfg.NewTracer()
}
