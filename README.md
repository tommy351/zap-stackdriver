# zap-stackdriver

[![GoDoc](https://godoc.org/github.com/tommy351/zap-stackdriver?status.svg)](https://godoc.org/github.com/tommy351/zap-stackdriver) [![Build Status](https://travis-ci.org/tommy351/zap-stackdriver.svg)](https://travis-ci.org/tommy351/zap-stackdriver) [![Coverage Status](https://coveralls.io/repos/github/tommy351/zap-stackdriver/badge.svg)](https://coveralls.io/github/tommy351/zap-stackdriver)

Prints [Stackdriver format](https://cloud.google.com/error-reporting/docs/formatting-error-messages) logs with [zap](https://github.com/uber-go/zap).

## Installation

``` sh
go get -u github.com/tommy351/zap-stackdriver
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
