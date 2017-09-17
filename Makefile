%_test.go:
	mockery -dir="vendor/go.uber.org/zap/zapcore" -name=$* -output="." -testonly -outpkg="stackdriver"

mocks: ObjectEncoder_test.go PrimitiveArrayEncoder_test.go
.PHONY: mocks
