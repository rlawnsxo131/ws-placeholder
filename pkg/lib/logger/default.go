package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	_onceDefaultLogger      sync.Once
	_singletonDefaultLogger *defaultLogger
)

func Default() *defaultLogger {
	_onceDefaultLogger.Do(func() {
		_singletonDefaultLogger = NewDefaultLogger(os.Stdout)
	})
	return _singletonDefaultLogger
}

type defaultLogger struct {
	*zerolog.Logger
}

func NewDefaultLogger(w io.Writer) *defaultLogger {
	l := zerolog.New(w).With().Caller().Timestamp().Logger()
	return &defaultLogger{&l}
}
