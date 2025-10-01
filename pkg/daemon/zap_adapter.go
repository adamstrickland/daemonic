package daemon

import "go.uber.org/zap"

type ZapAdapter struct {
	logger *zap.Logger
}

var _ Logger = (*ZapAdapter)(nil)

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{logger: logger}
}

func (z *ZapAdapter) Debug(msg string, args ...any) {
	z.logger.Sugar().Debug(msg, toZapFields(args))
}

func (z *ZapAdapter) Info(msg string, args ...any) {
	z.logger.Sugar().Info(msg, toZapFields(args))
}

func (z *ZapAdapter) Warn(msg string, args ...any) {
	z.logger.Sugar().Warn(msg, toZapFields(args))
}

func (z *ZapAdapter) Error(msg string, args ...any) {
	z.logger.Sugar().Error(msg, toZapFields(args))
}

func toZapFields(args []any) []zap.Field {
	fields := make([]zap.Field, len(args))
	for i, arg := range args {
		fields[i] = zap.Any("", arg)
	}
	return fields
}
