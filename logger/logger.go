package logger

type Logger interface {
	Log(Level Level, kvs ...interface{})
}

type LogContext struct {
	loggers []Logger
	prefix  []Valuer
}

func (lc LogContext) Log(Level Level, kvs ...interface{}) {
	vs := make([]interface{}, 0, len(lc.prefix)+len(kvs))
	for _, f := range lc.prefix {
		k, v := f()
		vs = append(vs, k, v)
	}
	vs = append(vs, kvs...)
	for _, l := range lc.loggers {
		l.Log(Level, vs...)
	}
}

// decorator

// add prefix info to logger
func WithPrefix(logger Logger, prefix ...Valuer) Logger {
	if lc, ok := logger.(LogContext); ok {
		vs := make([]Valuer, 0, len(lc.prefix)+len(prefix))
		vs = append(vs, prefix...)
		vs = append(vs, lc.prefix...)
		return &LogContext{
			loggers: lc.loggers,
			prefix:  vs,
		}
	}
	// init
	return &LogContext{loggers: []Logger{logger}, prefix: prefix}
}

// add logger in log chain tail
func WithLogger(oldLogger Logger, newLogs ...Logger) Logger {
	if lc, ok := oldLogger.(LogContext); ok {
		ls := make([]Logger, 0, len(lc.loggers)+len(newLogs))
		ls = append(ls, lc.loggers...)
		ls = append(ls, newLogs...)
		return &LogContext{
			loggers: ls,
			prefix:  lc.prefix,
		}
	}
	// init
	ls := make([]Logger, 0, len(newLogs)+1)
	ls = append(ls, oldLogger)
	ls = append(ls, newLogs...)
	return &LogContext{loggers: ls}
}
