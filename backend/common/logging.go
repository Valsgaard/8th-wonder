package common

import (
	"context"
	"io"
	"time"

	"github.com/Sirupsen/logrus"
)

type contextKey int

const (
	contextLog contextKey = iota
)

// Log retrieves a log from the given comtext
func Log(ctx context.Context) *logrus.Entry {
	return ctx.Value(contextLog).(*logrus.Entry)
}

// SetLog stores a logger in the given context
func SetLog(ctx context.Context, v *logrus.Entry) context.Context {
	return context.WithValue(ctx, contextLog, v)
}

// NewLog creates a new logger, writing to the given writer
func NewLog(service string, out io.Writer) *logrus.Entry {
	log := logrus.New()
	log.Out = out
	log.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}

	return log.WithField("service", service)
}
