// https://github.com/uber-go/fx/blob/v1.20.1/fxevent/zap.go

package logging

import (
	"context"
	"log/slog"
	"strings"

	"go.uber.org/fx/fxevent"
)

type fxeventLogger struct {
	logger     Logger
	eventLevel Level // default: LevelInfo
	errorLevel Level
}

var _ fxevent.Logger = (*fxeventLogger)(nil)

// NewFxeventLogger returns a new fxevent.Logger that logs to logger.
func NewFxeventLogger(logger Logger) *fxeventLogger {
	l := &fxeventLogger{
		logger:     logger,
		eventLevel: LevelInfo,
		errorLevel: LevelError,
	}
	return l
}

// UseEventLevel sets the level of non-error logs emitted by Fx to level.
func (l *fxeventLogger) UseEventLevel(level Level) *fxeventLogger {
	l.eventLevel = level
	return l
}

// UseErrorLevel sets the level of error logs emitted by Fx to level.
func (l *fxeventLogger) UseErrorLevel(level Level) *fxeventLogger {
	l.errorLevel = level
	return l
}

func (l *fxeventLogger) logEvent(msg string, args ...any) {
	l.logger.Log(context.Background(), l.eventLevel, msg, args...)
}

func (l *fxeventLogger) logError(msg string, args ...any) {
	l.logger.Log(context.Background(), l.errorLevel, msg, args...)
}

// LogEvent implements fxevent.Logger.
func (l *fxeventLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logEvent("OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logEvent("OnStart hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"error", e.Err,
			)
		} else {
			l.logEvent("OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		l.logEvent("OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logError("OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"error", e.Err,
			)
		} else {
			l.logEvent("OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.logError("error encountered while applying options",
				"type", e.TypeName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"error", e.Err)
		} else {
			l.logEvent("supplied",
				"type", e.TypeName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("provided",
				"constructor", e.ConstructorName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"type", rtype,
				maybeBool("private", e.Private),
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				moduleField(e.ModuleName),
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				"error", e.Err,
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("replaced",
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.logError("error encountered while replacing",
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"error", e.Err,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.logEvent("decorated",
				"decorator", e.DecoratorName,
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.logError("error encountered while applying options",
				"stacktrace", e.StackTrace,
				"moduletrace", e.ModuleTrace,
				moduleField(e.ModuleName),
				"error", e.Err,
			)
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.logError("error returned",
				"name", e.Name,
				"kind", e.Kind,
				moduleField(e.ModuleName),
				"error", e.Err,
			)
		} else {
			l.logEvent("run",
				"name", e.Name,
				"kind", e.Kind,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Invoking:
		l.logEvent("invoking",
			"function", e.FunctionName,
			moduleField(e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.logError("invoke failed",
				"error", e.Err,
				"stack", e.Trace,
				"function", e.FunctionName,
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Stopping:
		l.logEvent("received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logError("stop failed",
				"error", e.Err,
			)
		}
	case *fxevent.RollingBack:
		l.logError("start failed, rolling back",
			"error", e.StartErr,
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logError("rollback failed",
				"error", e.Err,
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logError("start failed",
				"error", e.Err,
			)
		} else {
			l.logEvent("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logError("custom logger initialization failed",
				"error", e.Err,
			)

		} else {
			l.logEvent("initialized custom Logger",
				"function", e.ConstructorName,
			)
		}
	}
}

// empty group will be optimized away
var slogAttrSkip = slog.Group("skipped")

func moduleField(name string) slog.Attr {
	if len(name) == 0 {
		return slogAttrSkip
	}
	return slog.String("module", name)
}

func maybeBool(name string, b bool) slog.Attr {
	if b {
		return slog.Bool(name, true)
	}
	return slogAttrSkip
}
