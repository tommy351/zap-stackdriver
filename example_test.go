package stackdriver_test

import (
	stackdriver "github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Example_basic() {
	config := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{
			Core: core,
		}
	}), zap.Fields(
		stackdriver.LogServiceContext(&stackdriver.ServiceContext{
			Service: "foo",
			Version: "bar",
		}),
	))

	if err != nil {
		panic(err)
	}

	logger.Info("Hello",
		stackdriver.LogUser("token"),
		stackdriver.LogHTTPRequest(&stackdriver.HTTPRequest{
			Method:             "GET",
			URL:                "/foo",
			UserAgent:          "bar",
			Referrer:           "baz",
			ResponseStatusCode: 200,
			RemoteIP:           "1.2.3.4",
		}))
}
