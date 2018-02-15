package stackdriver

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logKeyServiceContext        = "serviceContext"
	logKeyContextHTTPRequest    = "context.httpRequest"
	logKeyContextUser           = "context.user"
	logKeyContextReportLocation = "context.reportLocation"
)

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

var EncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "eventTime",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    EncodeLevel,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

type Core struct {
	zapcore.Core

	SetReportLocation bool

	ctx *Context
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	fields, ctx := c.extractCtx(fields)

	return &Core{
		Core:              c.Core.With(fields),
		SetReportLocation: c.SetReportLocation,
		ctx:               ctx,
	}
}

func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	loc := c.getReportLocationFromEntry(entry)

	if loc != nil {
		fields = append(fields, LogReportLocation(loc))
	}

	fields, ctx := c.extractCtx(fields)
	fields = append(fields, zap.Object("context", ctx))

	return c.Core.Write(entry, fields)
}

func (c *Core) extractCtx(fields []zapcore.Field) ([]zapcore.Field, *Context) {
	output := []zapcore.Field{}
	ctx := c.cloneCtx()

	for _, f := range fields {
		switch f.Key {
		case logKeyContextHTTPRequest:
			ctx.HTTPRequest = f.Interface.(*HTTPRequest)
		case logKeyContextReportLocation:
			ctx.ReportLocation = f.Interface.(*ReportLocation)
		case logKeyContextUser:
			ctx.User = f.String
		default:
			output = append(output, f)
		}
	}

	return output, ctx
}

func (c *Core) cloneCtx() *Context {
	if c.ctx == nil {
		return &Context{}
	}

	return c.ctx.Clone()
}

func (c *Core) getReportLocationFromEntry(entry zapcore.Entry) *ReportLocation {
	if !c.SetReportLocation {
		return nil
	}

	caller := entry.Caller

	if !caller.Defined {
		return nil
	}

	loc := &ReportLocation{
		FilePath:   caller.File,
		LineNumber: caller.Line,
	}

	if fn := runtime.FuncForPC(caller.PC); fn != nil {
		loc.FunctionName = fn.Name()
	}

	return loc
}

func LogServiceContext(ctx *ServiceContext) zapcore.Field {
	return zap.Object(logKeyServiceContext, ctx)
}

func LogHTTPRequest(req *HTTPRequest) zapcore.Field {
	return zap.Object(logKeyContextHTTPRequest, req)
}

func LogUser(user string) zapcore.Field {
	return zap.String(logKeyContextUser, user)
}

func LogReportLocation(loc *ReportLocation) zapcore.Field {
	return zap.Object(logKeyContextReportLocation, loc)
}

func EncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[lv])
}
