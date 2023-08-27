package logger

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

// FxLogger is an Fx event logger that logs events to Zero.
type FxLogger struct {
	Logger *zerolog.Logger
}

func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Trace().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.Logger.Trace().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	case *fxevent.OnStopExecuting:
		l.Logger.Trace().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).
				Str("callee", e.FunctionName).
				Str("callee", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.Logger.Trace().Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Logger.Warn().Err(e.Err).Str("type", e.TypeName).Msg("supplied")
		} else {
			l.Logger.Trace().Str("type", e.TypeName).Msg("supplied")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Logger.Trace().Str("type", rtype).
				Str("constructor", e.ConstructorName).
				Msg("Provided")
		}
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("error encountered while applying options")
		}
	case *fxevent.Invoking:
		// Do nothing. Will log on Invoked.

	case *fxevent.Invoked:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Str("stack", e.Trace).
				Str("function", e.FunctionName).Msg("invoke failed")
		} else {
			l.Logger.Trace().Str("function", e.FunctionName).Msg("invoked")
		}
	case *fxevent.Stopping:
		l.Logger.Info().Str("signal", strings.ToUpper(e.Signal.String())).Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.Logger.Error().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("start failed")
		} else {
			l.Logger.Trace().Msg("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Logger.Error().Err(e.Err).Msg("custom logger initialization failed")
		} else {
			l.Logger.Trace().Str("function", e.ConstructorName).Msg("Initialized custom fxevent.Logger")
		}
	}
}
