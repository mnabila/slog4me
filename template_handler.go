package slog4me

import (
	"context"
	"io"
	"log/slog"
	"slices"
	"sync"
	"text/template"

	"golang.org/x/sync/errgroup"
)

type Mapper[T any] func(record slog.Record) T

type TemplateWriter[T any] struct {
	Writer   io.Writer
	Template *template.Template
	Levels   []slog.Level
}

type TemplateHandler[T any] struct {
	mu      *sync.Mutex
	writers []TemplateWriter[T]
	mapper  Mapper[T]
}

func NewTemplateHandler[T any](opts ...HandlerOption[T]) (slog.Handler, error) {
	handler := &TemplateHandler[T]{
		mu:      &sync.Mutex{},
		writers: []TemplateWriter[T]{},
		mapper:  nil,
	}

	for _, opt := range opts {
		if err := opt(handler); err != nil {
			return nil, err
		}
	}

	return handler, nil
}

func (h *TemplateHandler[T]) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *TemplateHandler[T]) Handle(ctx context.Context, record slog.Record) error {
	data := h.mapper(record)

	var g errgroup.Group

	for _, writer := range h.writers {
		if slices.Contains(writer.Levels, record.Level) {
			w := writer
			g.Go(func() error {
				h.mu.Lock()
				defer h.mu.Unlock()
				return w.Template.Execute(w.Writer, data)
			})
		}
	}

	return g.Wait()
}

func (h *TemplateHandler[T]) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *TemplateHandler[T]) WithGroup(name string) slog.Handler {
	return h
}
