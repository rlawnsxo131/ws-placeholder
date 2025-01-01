package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	_onceDefaultLogger      sync.Once
	_singletonDefaultLogger *DefaultLogger
)

func Default() *DefaultLogger {
	_onceDefaultLogger.Do(func() {
		_singletonDefaultLogger = NewDefaultLogger(os.Stdout)
	})
	return _singletonDefaultLogger
}

type DefaultLogger struct {
	*zerolog.Logger
}

func NewDefaultLogger(w io.Writer) *DefaultLogger {
	l := zerolog.New(w).With().Caller().Timestamp().Logger()
	return &DefaultLogger{&l}
}
