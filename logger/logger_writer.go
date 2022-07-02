package logger

import (
	"github.com/rs/zerolog"
)

// Info level logger writer
type InfoLoggerWriter struct {
	w zerolog.LevelWriter
}

func (lw *InfoLoggerWriter) Write(p []byte) (n int, err error) {
	return lw.w.Write(p)
}

func (lw *InfoLoggerWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= zerolog.ErrorLevel {
		return lw.w.WriteLevel(level, p)
	}
	return len(p), nil
}

// Error level logger writer
type ErrorLoggerWriter struct {
	w zerolog.LevelWriter
}

func (lw *ErrorLoggerWriter) Write(p []byte) (n int, err error) {
	return lw.w.Write(p)
}

func (lw *ErrorLoggerWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < zerolog.ErrorLevel {
		return lw.w.WriteLevel(level, p)
	}
	return len(p), nil
}
