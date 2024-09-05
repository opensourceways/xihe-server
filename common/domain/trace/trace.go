// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Example using OTLP exporters + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type Trace interface {
	Span(ctx context.Context, name string) (spanctx context.Context, span trace.Span)
}
