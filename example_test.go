package ddawslogrus_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/nolotz/ddawslogrus"
)

func TestExample(t *testing.T) {
	rootContext := context.Background()

	// inject request & trace id's
	logrus.AddHook(
		ddawslogrus.NewHook().WithContextFunc(func() context.Context {
			return rootContext
		}),
	)

	// use common timestamp field
	logrus.SetFormatter(ddawslogrus.NewFormatter())

	// ...

	ctx := lambdacontext.NewContext(rootContext, new(lambdacontext.LambdaContext))
	ctx = tracer.ContextWithSpan(ctx, tracer.Span(nil))

	// ...

	logrus.WithContext(ctx).Info("enjoy bro")
}
