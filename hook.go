package ddawslogrus

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
)

func NewHook() *Hook {
	return &Hook{
		levels: logrus.AllLevels,
	}
}

type ContextFunc func() context.Context

type Hook struct {
	ctxFn  ContextFunc
	levels []logrus.Level
}

func (h *Hook) WithContextFunc(ctxFn ContextFunc) *Hook {
	h.ctxFn = ctxFn
	return h
}

func (h *Hook) WithLevels(levels ...logrus.Level) *Hook {
	h.levels = levels
	return h
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	var ctx context.Context
	if h.ctxFn != nil {
		ctx = h.ctxFn()
	}
	if entry.Context != nil {
		ctx = entry.Context
	}

	if ctx == nil {
		return nil
	}

	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if ok {
		if entry.Data["lambda"] == nil {
			entry.Data["lambda"] = make(map[string]any)
		}

		entry.Data["lambda"].(map[string]any)["request_id"] = lambdaCtx.AwsRequestID
	}

	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		entry.Data["dd.trace_id"] = span.Context().TraceID()
		entry.Data["dd.span_id"] = span.Context().SpanID()
	}

	return nil
}

func (h *Hook) Levels() []logrus.Level {
	return h.levels
}
