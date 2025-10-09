package bundle

import (
	"context"
	"xo/src/telemetry"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"go.opentelemetry.io/otel"

	oras "oras.land/oras-go/v2"
)

func Pull(ctx context.Context, repo oras.Target, target oras.Target, tag string) (v1.Descriptor, error) {
	_, span := otel.Tracer("xo").Start(ctx, "BundlePullV1")
	telemetry.SetSpanAttributes(span)
	defer span.End()

	return oras.Copy(ctx, repo, tag, target, tag, oras.DefaultCopyOptions)
}
