// Package logging configures the process-wide slog logger: output goes to
// stderr (visible in `wails dev`) and to a log file, because a packaged
// Windows build has no console.
package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	appDir      = "lol-teammate-helper"
	fileName    = "app.log"
	maxFileSize = 5 << 20 // restart the file once it grows past 5 MiB
	levelEnvVar = "LTH_LOG_LEVEL"
)

// Dir returns the directory holding the log file.
func Dir() string {
	base, err := os.UserConfigDir()
	if err != nil {
		base = os.TempDir()
	}
	return filepath.Join(base, appDir, "logs")
}

// Setup installs the default logger and returns the log file path (empty if
// the file could not be opened) and a function that closes the file.
// Set LTH_LOG_LEVEL=debug to see every LCU request.
func Setup() (path string, closeFn func()) {
	level := parseLevel(os.Getenv(levelEnvVar))

	var out io.Writer = os.Stderr
	closeFn = func() {}

	dir := Dir()
	if err := os.MkdirAll(dir, 0o755); err == nil {
		path = filepath.Join(dir, fileName)
		flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
		if info, err := os.Stat(path); err == nil && info.Size() > maxFileSize {
			flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
		}
		if f, err := os.OpenFile(path, flags, 0o644); err == nil {
			out = io.MultiWriter(os.Stderr, f)
			closeFn = func() { _ = f.Close() }
		} else {
			path = ""
		}
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: level})))
	slog.Info("logging started", "file", path, "level", level.String())
	return path, closeFn
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
