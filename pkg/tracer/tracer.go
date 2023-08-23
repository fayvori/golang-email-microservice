package tracer

import (
	"fmt"
	"go-email/config"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	service = "email-service"
)

// TracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func TracerProvider(c *config.Config) *tracesdk.TracerProvider {
	url := fmt.Sprintf("http://%s/api/traces", c.Jaeger.Host)

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		log.Errorf("Cannot connect to the jaeger %s", err.Error())
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(service),
		)),
	)

	return tp
}
