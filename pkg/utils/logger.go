package utils

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogField represents a logging field that can be used with our logger
type LogField struct {
	Key   string
	Value interface{}
	Type  LogFieldType
}

// LogFieldType represents the type of a field
type LogFieldType int

const (
	LogFieldTypeString LogFieldType = iota
	LogFieldTypeInt
	LogFieldTypeFloat64
	LogFieldTypeBool
	LogFieldTypeDuration
	LogFieldTypeTime
	LogFieldTypeError
)

// LogField constructors
func String(key, value string) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeString}
}

func Int(key string, value int) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeInt}
}

func Float64(key string, value float64) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeFloat64}
}

func Bool(key string, value bool) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeBool}
}

func Duration(key string, value time.Duration) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeDuration}
}

func Time(key string, value time.Time) LogField {
	return LogField{Key: key, Value: value, Type: LogFieldTypeTime}
}

func ErrorField(err error) LogField {
	return LogField{Key: "error", Value: err, Type: LogFieldTypeError}
}

// convertFields converts our custom LogField type to zap.Field
func convertFields(fields ...LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		switch field.Type {
		case LogFieldTypeString:
			zapFields[i] = zap.String(field.Key, field.Value.(string))
		case LogFieldTypeInt:
			zapFields[i] = zap.Int(field.Key, field.Value.(int))
		case LogFieldTypeFloat64:
			zapFields[i] = zap.Float64(field.Key, field.Value.(float64))
		case LogFieldTypeBool:
			zapFields[i] = zap.Bool(field.Key, field.Value.(bool))
		case LogFieldTypeDuration:
			zapFields[i] = zap.Duration(field.Key, field.Value.(time.Duration))
		case LogFieldTypeTime:
			zapFields[i] = zap.Time(field.Key, field.Value.(time.Time))
		case LogFieldTypeError:
			zapFields[i] = zap.Error(field.Value.(error))
		}
	}
	return zapFields
}

var Logger *zap.Logger

func InitLogger(level, format string) error {
	var config zap.Config

	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	if format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level = zap.NewAtomicLevelAt(zapLevel)

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Logger)

	return nil
}

func GetLogger() *zap.Logger {
	if Logger == nil {
		Logger, _ = zap.NewDevelopment()
	}
	return Logger
}

func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

func Info(msg string, fields ...LogField) {
	GetLogger().Info(msg, convertFields(fields...)...)
}

func Debug(msg string, fields ...LogField) {
	GetLogger().Debug(msg, convertFields(fields...)...)
}

func Warn(msg string, fields ...LogField) {
	GetLogger().Warn(msg, convertFields(fields...)...)
}

func Error(msg string, fields ...LogField) {
	GetLogger().Error(msg, convertFields(fields...)...)
}

func Fatal(msg string, fields ...LogField) {
	GetLogger().Fatal(msg, convertFields(fields...)...)
}
