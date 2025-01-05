package logger

import (
	"io"
	"os"
	"sync"
	"time"

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
	// ko: +9 hours
	l := zerolog.New(w).With().Caller().Str("time", time.Now().UTC().Format(time.RFC3339Nano)).Logger()
	return &DefaultLogger{&l}
}
