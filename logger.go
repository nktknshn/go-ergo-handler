package goergohandler

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
}

type LoggerType struct {
	logger Logger
}

func NewLogger(logger Logger) *LoggerType {
	return &LoggerType{
		logger,
	}
}

// func (l *LoggerType) Attach(builder *Builder[]) {

// }

type AttacherLoggerType struct {
	l *LoggerType
}

func (a *AttacherLoggerType) Info(args ...any) {
	a.l.logger.Info(args...)
}
