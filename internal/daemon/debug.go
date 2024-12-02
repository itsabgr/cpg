package daemon

import (
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"os"
	"slices"
	"strings"
)

var debug = os.Getenv("DEBUG") != "" || slices.ContainsFunc(os.Args, func(s string) bool {
	return strings.Contains(s, "debug")
})

func Debug() bool {
	return debug
}

func init() {
	if debug {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
		})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		})))
	}
	slog.Debug("debug mode", "debug", debug)
}
