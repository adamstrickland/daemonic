package daemon

import (
	"log/slog"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type SlogAdapter struct {
	logger *slog.Logger
}

func NewSlogAdapter(logger *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: logger}
}

func (s *SlogAdapter) Debug(msg string, args ...any) {
	s.logger.Debug(msg, args...)
}

func (s *SlogAdapter) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}

func (s *SlogAdapter) Warn(msg string, args ...any) {
	s.logger.Warn(msg, args...)
}

func (s *SlogAdapter) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{logger: logger}
}

func (z *ZapAdapter) Debug(msg string, args ...any) {
	z.logger.Debug(msg, toZapFields(args)...)
}

func (z *ZapAdapter) Info(msg string, args ...any) {
	z.logger.Info(msg, toZapFields(args)...)
}

func (z *ZapAdapter) Warn(msg string, args ...any) {
	z.logger.Warn(msg, toZapFields(args)...)
}

func (z *ZapAdapter) Error(msg string, args ...any) {
	z.logger.Error(msg, toZapFields(args)...)
}

func toZapFields(args []any) []zap.Field {
	fields := make([]zap.Field, len(args))
	for i, arg := range args {
		fields[i] = zap.Any("", arg)
	}
	return fields
}
