package stackdriver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logEntry struct {
	Severity       string          `json:"severity"`
	EventTime      logEntryTime    `json:"eventTime"`
	ServiceContext *ServiceContext `json:"serviceContext"`
	Message        string          `json:"message"`
	Context        *Context        `json:"context"`
}

type logEntryTime time.Time

func (t *logEntryTime) UnmarshalText(text []byte) error {
	res, err := time.Parse("2006-01-02T15:04:05.000Z0700", string(text))

	if err != nil {
		return err
	}

	*t = logEntryTime(res)
	return nil
}

func newCore(writer io.Writer) *Core {
	enc := zapcore.NewJSONEncoder(EncoderConfig)
	core := zapcore.NewCore(enc, zapcore.AddSync(writer), zapcore.DebugLevel)

	return &Core{
		Core: core,
	}
}

func TestCore(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	core := newCore(writer)
	logger := zap.New(core)

	t.Run("Basic", func(t *testing.T) {
		defer writer.Reset()

		subLogger := logger.With(LogServiceContext(&ServiceContext{
			Service: "foo",
			Version: "bar",
		}))
		subLogger.Debug("test", zap.String("foo", "bar"))

		var actual struct {
			logEntry

			Foo string `json:"foo"`
		}

		require.Nil(t, json.Unmarshal(writer.Bytes(), &actual))
		assert.Equal(t, "DEBUG", actual.Severity)
		assert.Equal(t, "test", actual.Message)
		assert.WithinDuration(t, time.Now(), time.Time(actual.EventTime), time.Second)
		assert.Equal(t, &ServiceContext{
			Service: "foo",
			Version: "bar",
		}, actual.ServiceContext)
		assert.Equal(t, "bar", actual.Foo)
	})

	t.Run("With context", func(t *testing.T) {
		defer writer.Reset()

		req := &HTTPRequest{
			Method:             "GET",
			URL:                "/foo",
			UserAgent:          "bar",
			Referrer:           "baz",
			ResponseStatusCode: 200,
			RemoteIP:           "1.2.3.4",
		}

		loc := &ReportLocation{
			FilePath:     "foo",
			FunctionName: "bar",
			LineNumber:   42,
		}

		logger.Debug("",
			LogUser("foo"),
			LogHTTPRequest(req),
			LogReportLocation(loc),
		)

		var actual logEntry
		require.Nil(t, json.Unmarshal(writer.Bytes(), &actual))
		assert.Equal(t, &Context{
			User:           "foo",
			HTTPRequest:    req,
			ReportLocation: loc,
		}, actual.Context)
	})

	t.Run("Set report location from entry", func(t *testing.T) {
		defer writer.Reset()

		core := newCore(writer)
		core.SetReportLocation = true
		logger := zap.New(core, zap.AddCaller())
		// nolint: goerr113
		err := errors.New("random error")
		_, file, line, _ := runtime.Caller(0)
		logger.Error("", zap.Error(err))

		var actual logEntry
		require.Nil(t, json.Unmarshal(writer.Bytes(), &actual))
		loc := actual.Context.ReportLocation
		assert.Equal(t, file, loc.FilePath)
		assert.Equal(t, line+1, loc.LineNumber)
		assert.True(t, strings.HasPrefix(loc.FunctionName, "github.com/tommy351/zap-stackdriver.TestCore"))
	})
}

func TestLogServiceContext(t *testing.T) {
	ctx := &ServiceContext{}
	field := LogServiceContext(ctx)
	assert.Equal(t, zap.Object(logKeyServiceContext, ctx), field)
}

func TestLogHTTPRequest(t *testing.T) {
	req := &HTTPRequest{}
	field := LogHTTPRequest(req)
	assert.Equal(t, zap.Object(logKeyContextHTTPRequest, req), field)
}

func TestLogUser(t *testing.T) {
	field := LogUser("foo")
	assert.Equal(t, zap.String(logKeyContextUser, "foo"), field)
}

func TestLogReportLocation(t *testing.T) {
	loc := &ReportLocation{}
	field := LogReportLocation(loc)
	assert.Equal(t, zap.Object(logKeyContextReportLocation, loc), field)
}

func TestEncodeLevel(t *testing.T) {
	tests := []struct {
		Level    zapcore.Level
		Expected string
	}{
		{
			Level:    zapcore.DebugLevel,
			Expected: "DEBUG",
		},
		{
			Level:    zapcore.InfoLevel,
			Expected: "INFO",
		},
		{
			Level:    zapcore.WarnLevel,
			Expected: "WARNING",
		},
		{
			Level:    zapcore.ErrorLevel,
			Expected: "ERROR",
		},
		{
			Level:    zapcore.DPanicLevel,
			Expected: "CRITICAL",
		},
		{
			Level:    zapcore.PanicLevel,
			Expected: "ALERT",
		},
		{
			Level:    zapcore.FatalLevel,
			Expected: "EMERGENCY",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.Expected, func(t *testing.T) {
			enc := new(PrimitiveArrayEncoder)
			enc.On("AppendString", test.Expected).Once()
			EncodeLevel(test.Level, enc)
			enc.AssertExpectations(t)
		})
	}
}
