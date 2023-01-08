package ddawslogrus_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/nolotz/ddawslogrus"
)

func TestHookFireLambda(t *testing.T) {
	ctx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: "1",
	})

	hook := ddawslogrus.NewHook().WithContext(ctx)
	entry := logrus.NewEntry(nil).WithFields(logrus.Fields{
		"lambda": map[string]any{
			"arn": "test",
		},
	})

	err := hook.Fire(entry)

	assert.NoError(t, err)
	assert.Equal(t, "1", entry.Data["lambda"].(map[string]any)["request_id"])
	assert.Equal(t, "test", entry.Data["lambda"].(map[string]any)["arn"])
}

func TestHookFireTracer(t *testing.T) {
	mockTracer := mocktracer.Start().(ddtrace.Tracer)
	defer mockTracer.Stop()

	span := mockTracer.StartSpan("test")
	ctx := tracer.ContextWithSpan(context.Background(), span)

	hook := ddawslogrus.NewHook()
	entry := logrus.NewEntry(nil).WithContext(ctx)

	err := hook.Fire(entry)

	assert.NoError(t, err)
	assert.Equal(t, span.Context().TraceID(), entry.Data["dd.trace_id"])
	assert.Equal(t, span.Context().SpanID(), entry.Data["dd.span_id"])
}

func TestHookLevels(t *testing.T) {
	hook := ddawslogrus.NewHook().WithLevels(logrus.DebugLevel)

	assert.Equal(t, []logrus.Level{logrus.DebugLevel}, hook.Levels())
}

func TestHookNoContext(t *testing.T) {
	hook := ddawslogrus.NewHook()
	entry := logrus.NewEntry(nil)

	err := hook.Fire(entry)

	assert.NoError(t, err)
}
