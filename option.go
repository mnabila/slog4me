package slog4me

import (
	"io"
	"log/slog"
	"text/template"
)

type HandlerOption[T any] func(*TemplateHandler[T]) error

func WithTemplateWriter[T any](w io.Writer, layout string, levels ...slog.Level) HandlerOption[T] {
	return func(th *TemplateHandler[T]) error {
		tmpl, err := template.New("_").Parse(layout)
		if err != nil {
			return err
		}

		th.writers = append(th.writers, TemplateWriter[T]{
			Writer:   w,
			Template: tmpl,
			Levels:   levels,
		})

		return nil
	}
}

func WithMapper[T any](mapper Mapper[T]) HandlerOption[T] {
	return func(opts *TemplateHandler[T]) error {
		opts.mapper = mapper
		return nil
	}
}
