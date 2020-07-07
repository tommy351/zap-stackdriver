# zap-stackdriver

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/tommy351/zap-stackdriver)](https://github.com/tommy351/zap-stackdriver/releases) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/tommy351/zap-stackdriver) ![Test](https://github.com/tommy351/zap-stackdriver/workflows/Test/badge.svg) [![codecov](https://codecov.io/gh/tommy351/zap-stackdriver/branch/master/graph/badge.svg)](https://codecov.io/gh/tommy351/zap-stackdriver)

Prints [Stackdriver format](https://cloud.google.com/error-reporting/docs/formatting-error-messages) logs with [zap](https://github.com/uber-go/zap).

## Installation

``` sh
go get github.com/tommy351/zap-stackdriver
```

## Usage

``` go
package main

import (
	"github.com/tommy351/zap-stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
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
```
