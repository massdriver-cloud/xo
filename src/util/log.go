package util

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func LogError(err error, span trace.Span, msg string) {
	log.Error().Err(err).Msg(msg)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
