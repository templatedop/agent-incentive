package log

import (
	"context"

	"github.com/rs/zerolog"
)

const (
	level   = "level"
	message = "message"
	service = "service"
	version = "version"
	stdout  = "stdout"
	noop    = "noop"
	test    = "test"
	console = "console"
	tagsKey = "tags" // Key for tags in log entries
)

type contextKey string

const (
	logTagsContextKey contextKey = "log_tags"
)

var (
	baseLogger     *Logger
	samplingConfig *SamplingConfig
)

// Tunable skip frame counts for caller attribution. These were reduced by 1
// after verification tests showed the previous values pointed to wrapper
// helpers instead of the user's call site.
// - eventAPISkipFrames: when using the Event API (InfoEvent, DebugEvent, ...)
// - simpleAPISkipFrames: when using the convenience message API (Info, Debug, ...)
// If additional wrappers are introduced, adjust these constants or make them
// configurable.
const (
	eventAPISkipFrames  = 1 // previously 2
	simpleAPISkipFrames = 2 // previously 3
)

type Logger struct {
	logger *zerolog.Logger
}

// ToZerolog exposes the internal zerolog logger.
// Deprecated: Use the Event-based API instead (DebugEvent, InfoEvent, etc.)
// for structured logging with fields. This method will be removed in a future version.
//
// Instead of:
//
//	log.GetBaseLoggerInstance().ToZerolog().Info().Str("key", "val").Msg("message")
//
// Use:
//
//	log.InfoEvent(ctx).Str("key", "val").Msg("message")
func (l *Logger) ToZerolog() *zerolog.Logger {
	return l.logger
}

// FromZerolog converts a zerolog logger to a logger.
// Deprecated: This is mainly for internal use. Use the standard API instead.
func FromZerolog(logger zerolog.Logger) *Logger {
	return &Logger{&logger}
}

// Debug logs a debug message. Supports string messages, errors, and format strings.
// For structured logging with fields, use DebugEvent instead.
func Debug(ctx context.Context, message interface{}, args ...interface{}) {
	logWithEvent(getEventLoggerWithSkip(ctx, zerolog.DebugLevel, simpleAPISkipFrames), message, args...)
}

// Info logs an info message. Supports string messages, errors, and format strings.
// For structured logging with fields, use InfoEvent instead.
func Info(ctx context.Context, message interface{}, args ...interface{}) {
	logWithEvent(getEventLoggerWithSkip(ctx, zerolog.InfoLevel, simpleAPISkipFrames), message, args...)
}

// Warn logs a warning message. Supports string messages, errors, and format strings.
// For structured logging with fields, use WarnEvent instead.
func Warn(ctx context.Context, message interface{}, args ...interface{}) {
	logWithEvent(getEventLoggerWithSkip(ctx, zerolog.WarnLevel, simpleAPISkipFrames), message, args...)
}

// Error logs an error message. Supports string messages, errors, and format strings.
// For structured logging with fields, use ErrorEvent instead.
func Error(ctx context.Context, message interface{}, args ...interface{}) {
	logWithEvent(getEventLoggerWithSkip(ctx, zerolog.ErrorLevel, simpleAPISkipFrames), message, args...)
}

// Critical logs a message at FatalLevel but does not exit the application.
// Use this for critical errors that should be logged at the highest severity
// but where you want to handle cleanup or continue execution.
// For structured logging with fields, use CriticalEvent instead.
func Critical(ctx context.Context, message interface{}, args ...interface{}) {
	logWithEvent(getEventLoggerWithSkip(ctx, zerolog.FatalLevel, simpleAPISkipFrames), message, args...)
}

// Fatal logs a message at FatalLevel and exits the application with status code 1.
// This follows standard logging library conventions where Fatal terminates the program.
// For critical errors that don't require immediate exit, use Critical instead.
func Fatal(ctx context.Context, message interface{}, args ...interface{}) {
	Critical(ctx, message, args...)
	panic("fatal error occurred") // Using panic instead of os.Exit for better testability
}

// DebugWithFields logs a debug message with structured fields.
// This is a convenience function that combines simple message logging with structured fields.
//
// Example:
//
//	log.DebugWithFields(ctx, "processing user", map[string]interface{}{
//	    "user_id": "123",
//	    "action": "login",
//	})

func DebugWithFields(ctx context.Context, message string, fields map[string]interface{}) {
	event := getEventLoggerWithSkip(ctx, zerolog.DebugLevel, eventAPISkipFrames)
	addFieldsToEvent(event, fields)
	event.Msg(message)
}

// InfoWithFields logs an info message with structured fields.
// This is a convenience function that combines simple message logging with structured fields.
//
// Example:
//
//	log.InfoWithFields(ctx, "user logged in", map[string]interface{}{
//	    "user_id": "123",
//	    "ip": "192.168.1.1",
//	})
func InfoWithFields(ctx context.Context, message string, fields map[string]interface{}) {
	event := getEventLoggerWithSkip(ctx, zerolog.InfoLevel, eventAPISkipFrames)
	addFieldsToEvent(event, fields)
	event.Msg(message)
}

// WarnWithFields logs a warning message with structured fields.
// This is a convenience function that combines simple message logging with structured fields.
//
// Example:
//
//	log.WarnWithFields(ctx, "rate limit approaching", map[string]interface{}{
//	    "attempts": 4,
//	    "limit": 5,
//	})
func WarnWithFields(ctx context.Context, message string, fields map[string]interface{}) {
	event := getEventLoggerWithSkip(ctx, zerolog.WarnLevel, eventAPISkipFrames)
	addFieldsToEvent(event, fields)
	event.Msg(message)
}

// ErrorWithFields logs an error message with structured fields.
// This is a convenience function that combines simple message logging with structured fields.
//
// Example:
//
//	log.ErrorWithFields(ctx, "database query failed", map[string]interface{}{
//	    "error": err,
//	    "query": sql,
//	    "duration": elapsed,
//	})
func ErrorWithFields(ctx context.Context, message string, fields map[string]interface{}) {
	event := getEventLoggerWithSkip(ctx, zerolog.ErrorLevel, eventAPISkipFrames)
	addFieldsToEvent(event, fields)
	event.Msg(message)
}

// CriticalWithFields logs a critical message with structured fields.
// This is a convenience function that combines simple message logging with structured fields.
//
// Example:
//
//	log.CriticalWithFields(ctx, "service unavailable", map[string]interface{}{
//	    "service": "payment-gateway",
//	    "error": err,
//	})
func CriticalWithFields(ctx context.Context, message string, fields map[string]interface{}) {
	event := getEventLoggerWithSkip(ctx, zerolog.FatalLevel, eventAPISkipFrames)
	addFieldsToEvent(event, fields)
	event.Msg(message)
}

// DebugEvent returns a zerolog.Event for structured logging at Debug level.
// This allows adding fields before calling Msg() to log the event.
//
// Example:
//
//	log.DebugEvent(ctx).Str("user_id", "123").Int("count", 10).Msg("processing items")
func DebugEvent(ctx context.Context) *zerolog.Event {
	return getEventLoggerWithSkip(ctx, zerolog.DebugLevel, eventAPISkipFrames)
}

// InfoEvent returns a zerolog.Event for structured logging at Info level.
// This allows adding fields before calling Msg() to log the event.
//
// Example:
//
//	log.InfoEvent(ctx).Str("operation", "login").Dur("latency", duration).Msg("user logged in")
func InfoEvent(ctx context.Context) *zerolog.Event {
	return getEventLoggerWithSkip(ctx, zerolog.InfoLevel, eventAPISkipFrames)
}

// WarnEvent returns a zerolog.Event for structured logging at Warn level.
// This allows adding fields before calling Msg() to log the event.
//
// Example:
//
//	log.WarnEvent(ctx).Str("reason", "rate_limit").Int("attempts", 5).Msg("rate limit approaching")
func WarnEvent(ctx context.Context) *zerolog.Event {
	return getEventLoggerWithSkip(ctx, zerolog.WarnLevel, eventAPISkipFrames)
}

// ErrorEvent returns a zerolog.Event for structured logging at Error level.
// This allows adding fields before calling Msg() to log the event.
//
// Example:
//
//	log.ErrorEvent(ctx).Err(err).Str("query", sql).Msg("database query failed")
func ErrorEvent(ctx context.Context) *zerolog.Event {
	return getEventLoggerWithSkip(ctx, zerolog.ErrorLevel, eventAPISkipFrames)
}

// CriticalEvent returns a zerolog.Event for structured logging at Fatal level (without exiting).
// This allows adding fields before calling Msg() to log the event.
//
// Example:
//
//	log.CriticalEvent(ctx).Err(err).Str("service", "payment").Msg("payment service unavailable")
func CriticalEvent(ctx context.Context) *zerolog.Event {
	return getEventLoggerWithSkip(ctx, zerolog.FatalLevel, eventAPISkipFrames)
}

// getEventLoggerWithSkip returns a zerolog.Event with caller information for the given context and level.
// It handles context-aware logger retrieval and sets up proper caller skip frames.
// The skipFrames parameter controls how many stack frames to skip for accurate caller reporting.
func getEventLoggerWithSkip(ctx context.Context, level zerolog.Level, skipFrames int) *zerolog.Event {
	var logger *Logger
	if ctx != nil {
		logger = getCtxLogger(ctx)
	} else if baseLogger != nil {
		logger = baseLogger
	} else {
		logger = getDefaultLogger()
	}

	// Check sampling config to determine if this log should be emitted
	tags := GetTags(ctx)
	if samplingConfig != nil && !samplingConfig.ShouldLog(level, tags) {
		// Return a disabled event that will not log anything
		nopLogger := zerolog.Nop()
		return nopLogger.WithLevel(level)
	}

	// Add caller information with correct skip frame count.
	// Verified via tests:
	// - eventAPISkipFrames (1) for Event API usage: getEventLoggerWithSkip -> InfoEvent -> caller
	// - simpleAPISkipFrames (2) for simple API usage: getEventLoggerWithSkip -> Info -> logWithEvent -> caller
	// If this library adds/removes wrapper layers, adjust the constants above.
	lw := logger.logger.With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrames).Logger()
	event := lw.WithLevel(level)

	// Add tags from context if present
	if len(tags) > 0 {
		event.Strs(tagsKey, tags)
	}

	return event
}

// logWithEvent handles message formatting and logging via a zerolog.Event.
// This is used internally by the simple logging functions (Debug, Info, etc.)
// to ensure both simple and structured logging use the same code path.
func logWithEvent(event *zerolog.Event, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		// Log error's message
		if len(args) == 0 {
			event.Msg(msg.Error())
		} else {
			event.Msgf(msg.Error(), args...)
		}
	case string:
		// Log string message
		if len(args) == 0 {
			event.Msg(msg)
		} else {
			event.Msgf(msg, args...)
		}
	default:
		// Handle unexpected types
		event.Msgf("message %v has unknown type %T", message, message)
	}
}

// WithTags adds tags to the context. Tags are automatically included in all log entries
// created with this context.
//
// Example:
//
//	ctx = log.WithTags(ctx, "database", "payment")
//	log.Info(ctx, "processing transaction") // Will include tags: ["database", "payment"]
func WithTags(ctx context.Context, tags ...string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	existingTags := GetTags(ctx)
	// Pre-allocate slice to avoid growth
	allTags := make([]string, 0, len(existingTags)+len(tags))
	allTags = append(allTags, existingTags...)
	allTags = append(allTags, tags...)

	return context.WithValue(ctx, logTagsContextKey, allTags)
}

// GetTags retrieves tags from the context.
// Returns an empty slice if no tags are present.
func GetTags(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}

	if tags, ok := ctx.Value(logTagsContextKey).([]string); ok {
		return tags
	}

	return nil
}

// addFieldsToEvent adds all fields from a map to a zerolog.Event.
// This helper function handles type conversion for common Go types.
func addFieldsToEvent(event *zerolog.Event, fields map[string]interface{}) {
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event.Str(key, v)
		case int:
			event.Int(key, v)
		case int64:
			event.Int64(key, v)
		case int32:
			event.Int32(key, v)
		case int16:
			event.Int16(key, v)
		case int8:
			event.Int8(key, v)
		case uint:
			event.Uint(key, v)
		case uint64:
			event.Uint64(key, v)
		case uint32:
			event.Uint32(key, v)
		case uint16:
			event.Uint16(key, v)
		case uint8:
			event.Uint8(key, v)
		case float64:
			event.Float64(key, v)
		case float32:
			event.Float32(key, v)
		case bool:
			event.Bool(key, v)
		case error:
			event.Err(v)
		case []string:
			event.Strs(key, v)
		case []int:
			event.Ints(key, v)
		case []int64:
			event.Ints64(key, v)
		case nil:
			// Skip nil values
		default:
			// For complex types, use Interface which will JSON marshal
			event.Interface(key, v)
		}
	}
}
