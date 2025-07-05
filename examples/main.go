package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/mnabila/slog4me"
)

type Data struct {
	Timestamp string
	Level     string
	System    string
	Arch      string
	Message   string
	Fields    string
}

func main() {
	dfile, _ := os.Create("debug.log")
	ifile, _ := os.Create("info.log")

	dt := "{{.Timestamp}}|{{.System}}|{{.Arch}}|{{.Level}}|{{.Message}}|{{.Fields}}\n"
	it := "{{.Timestamp}}|{{.System}}|{{.Arch}}|{{.Level}}|{{.Message}}\n"

	handler, err := slog4me.NewTemplateHandler(
		slog4me.WithTemplateWriter[Data](dfile, dt, slog.LevelDebug, slog.LevelInfo, slog.LevelError),
		slog4me.WithTemplateWriter[Data](ifile, it, slog.LevelInfo, slog.LevelError),
		slog4me.WithMapper(func(record slog.Record) Data {
			data := Data{
				Timestamp: record.Time.Format("2006-01-02 03:04:05"),
				Level:     record.Level.String(),
				Arch:      "",
				Message:   record.Message,
				Fields:    "",
				System:    "",
			}

			var parts []string

			record.Attrs(func(a slog.Attr) bool {
				switch a.Key {
				case "System":
					data.System = a.Value.String()
				case "Arch":
					data.Arch = a.Value.String()
				default:
					parts = append(parts, a.Value.String())
				}
				return true
			})
			data.Fields = strings.Join(parts, ":")

			return data
		}),
	)
	if err != nil {
		panic(err)
	}

	logger := slog.New(handler)

	for i := range 100 {
		logger.Debug(
			fmt.Sprintf("debug %d", i),
			slog.String("System", runtime.GOOS),
			slog.String("Arch", runtime.GOARCH),
			slog.Any("User", map[string]any{"name": "robert", "age": 3 * i}),
		)

		logger.Info(
			fmt.Sprintf("info %d", i),
			slog.String("System", runtime.GOOS),
			slog.String("Arch", runtime.GOARCH),
		)

		logger.Error(
			fmt.Sprintf("loop %d", i),
			slog.String("System", runtime.GOOS),
			slog.String("Arch", runtime.GOARCH),
			slog.Any("User", map[string]any{"name": "robert", "age": 3 * i}),
		)
	}
}
