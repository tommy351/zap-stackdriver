package stackdriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceContext_Clone(t *testing.T) {
	src := &ServiceContext{
		Service: "foo",
		Version: "bar",
	}

	res := src.Clone()
	assert.Equal(t, src, res)
}

func TestServiceContext_MarshalLogObject(t *testing.T) {
	enc := new(ObjectEncoder)
	ctx := &ServiceContext{
		Service: "foo",
		Version: "bar",
	}

	enc.On("AddString", "service", ctx.Service).Once()
	enc.On("AddString", "version", ctx.Version).Once()
	require.Nil(t, ctx.MarshalLogObject(enc))
	enc.AssertExpectations(t)
}

func TestContext_Clone(t *testing.T) {
	src := &Context{
		User:           "foo",
		HTTPRequest:    &HTTPRequest{},
		ReportLocation: &ReportLocation{},
	}

	res := src.Clone()
	assert.Equal(t, src, res)
}

func TestContext_MarshalLogObject(t *testing.T) {
	enc := new(ObjectEncoder)
	ctx := &Context{
		User:           "foo",
		HTTPRequest:    &HTTPRequest{},
		ReportLocation: &ReportLocation{},
	}

	enc.On("AddString", "user", ctx.User).Once()
	enc.On("AddObject", "httpRequest", ctx.HTTPRequest).Return(nil).Once()
	enc.On("AddObject", "reportLocation", ctx.ReportLocation).Return(nil).Once()
	require.Nil(t, ctx.MarshalLogObject(enc))
	enc.AssertExpectations(t)
}

func TestHTTPRequest_Clone(t *testing.T) {
	src := &HTTPRequest{
		Method:             "GET",
		URL:                "/foo",
		UserAgent:          "bar",
		Referrer:           "baz",
		ResponseStatusCode: 200,
		RemoteIP:           "1.2.3.4",
	}

	res := src.Clone()
	assert.Equal(t, src, res)
}

func TestHTTPRequest_MarshalLogObject(t *testing.T) {
	enc := new(ObjectEncoder)
	req := &HTTPRequest{
		Method:             "GET",
		URL:                "/foo",
		UserAgent:          "bar",
		Referrer:           "baz",
		ResponseStatusCode: 200,
		RemoteIP:           "1.2.3.4",
	}

	enc.On("AddString", "method", req.Method).Once()
	enc.On("AddString", "url", req.URL).Once()
	enc.On("AddString", "userAgent", req.UserAgent).Once()
	enc.On("AddString", "referrer", req.Referrer).Once()
	enc.On("AddInt", "responseStatusCode", req.ResponseStatusCode).Once()
	enc.On("AddString", "remoteIp", req.RemoteIP).Once()
	require.Nil(t, req.MarshalLogObject(enc))
	enc.AssertExpectations(t)
}

func TestReportLocation_Clone(t *testing.T) {
	src := &ReportLocation{
		FilePath:     "foo",
		FunctionName: "bar",
		LineNumber:   42,
	}

	res := src.Clone()
	assert.Equal(t, src, res)
}

func TestReportLocation_MarshalLogObject(t *testing.T) {
	enc := new(ObjectEncoder)
	loc := &ReportLocation{
		FilePath:     "foo",
		FunctionName: "bar",
		LineNumber:   42,
	}

	enc.On("AddString", "filePath", loc.FilePath).Once()
	enc.On("AddString", "functionName", loc.FunctionName).Once()
	enc.On("AddInt", "lineNumber", loc.LineNumber).Once()
	require.Nil(t, loc.MarshalLogObject(enc))
	enc.AssertExpectations(t)
}
