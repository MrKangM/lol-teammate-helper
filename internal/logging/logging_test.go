package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := map[string]slog.Level{"": slog.LevelInfo, "debug": slog.LevelDebug, " DEBUG ": slog.LevelDebug, "warn": slog.LevelWarn, "error": slog.LevelError, "x": slog.LevelInfo}
	for in, want := range cases {
		if got := parseLevel(in); got != want {
			t.Errorf("parseLevel(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestSetupWritesFile(t *testing.T) {
	cfg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", cfg)
	t.Setenv("APPDATA", cfg)
	t.Setenv("HOME", cfg)

	path, closeFn := Setup()
	defer closeFn()
	slog.Info("hello from test")

	if path == "" || !strings.HasPrefix(path, cfg) {
		t.Fatalf("unexpected log path %q (config dir %q)", path, cfg)
	}
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil || !strings.Contains(string(data), "hello from test") {
		t.Fatalf("log file missing entry: %v %q", err, data)
	}
}
