// Package lcu locates a running League client and extracts its API credentials.
package lcu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Credentials are the per-launch connection details of the League client.
type Credentials struct {
	Port   int
	Token  string
	Region string
}

var (
	portRe   = regexp.MustCompile(`(?i)--app-port[=\s]+"?(\d+)`)
	tokenRe  = regexp.MustCompile(`(?i)--remoting-auth-token[=\s]+"?([\w-]+)`)
	regionRe = regexp.MustCompile(`(?i)--rso_platform_id[=\s]+"?([\w-]+)`)
)

// ErrNotRunning means no League client process was found.
var ErrNotRunning = errors.New("league client is not running")

// ErrAccessDenied means the client is running but its command line cannot be
// read, which happens when it runs elevated and this app does not.
var ErrAccessDenied = errors.New("league client found but its command line is unreadable: run this app as administrator, or set LTH_LOL_DIR to the folder containing the client's lockfile")

// Detect finds the running League client and returns its credentials.
// It reads the process command line first and falls back to the lockfile
// the client keeps in its install directory.
func Detect() (Credentials, error) {
	info, err := clientProcess()
	if err == nil {
		if creds, perr := ParseCommandLine(info.CommandLine); perr == nil {
			return creds, nil
		}
	}

	for _, dir := range lockfileDirs(info.InstallDir) {
		data, rerr := os.ReadFile(filepath.Join(dir, "lockfile"))
		if rerr != nil {
			continue
		}
		if creds, perr := ParseLockfile(string(data)); perr == nil {
			return creds, nil
		}
	}

	if err == nil {
		err = errors.New("league client command line has no credentials")
	}
	return Credentials{}, err
}

// ParseLockfile parses "LeagueClient:<pid>:<port>:<password>:<protocol>".
// The lockfile does not contain the region.
func ParseLockfile(text string) (Credentials, error) {
	parts := strings.Split(strings.TrimSpace(text), ":")
	if len(parts) < 5 {
		return Credentials{}, errors.New("malformed lockfile")
	}
	port, err := strconv.Atoi(parts[2])
	if err != nil || port <= 0 {
		return Credentials{}, errors.New("invalid port in lockfile")
	}
	if parts[3] == "" {
		return Credentials{}, errors.New("empty password in lockfile")
	}
	return Credentials{Port: port, Token: parts[3]}, nil
}

func lockfileDirs(installDir string) []string {
	var dirs []string
	if d := strings.TrimSpace(os.Getenv("LTH_LOL_DIR")); d != "" {
		dirs = append(dirs, d)
	}
	if installDir != "" {
		dirs = append(dirs, installDir)
	}
	return dirs
}

// ParseCommandLine extracts credentials from a LeagueClientUx command line.
func ParseCommandLine(text string) (Credentials, error) {
	var c Credentials

	m := portRe.FindStringSubmatch(text)
	if len(m) != 2 {
		return c, errors.New("app port not found")
	}
	port, err := strconv.Atoi(m[1])
	if err != nil {
		return c, fmt.Errorf("invalid port value: %w", err)
	}
	c.Port = port

	if m = tokenRe.FindStringSubmatch(text); len(m) != 2 {
		return c, errors.New("auth token not found")
	}
	c.Token = strings.TrimSpace(m[1])

	if m = regionRe.FindStringSubmatch(text); len(m) == 2 {
		c.Region = strings.TrimSpace(m[1])
	}
	return c, nil
}

type processInfo struct {
	CommandLine string
	InstallDir  string
}

func clientProcess() (processInfo, error) {
	if runtime.GOOS == "windows" {
		return windowsProcess()
	}
	out, err := exec.Command("ps", "-A", "-o", "args").Output()
	if err != nil {
		return processInfo{}, err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "LeagueClientUx") {
			return processInfo{CommandLine: line}, nil
		}
	}
	return processInfo{}, ErrNotRunning
}

// windowsProcess queries the process through PowerShell first and falls back
// to wmic, which is removed from recent Windows 11 builds.
func windowsProcess() (processInfo, error) {
	const script = `$p = Get-CimInstance Win32_Process -Filter "name='LeagueClientUx.exe'" | Select-Object -First 1
if (-not $p) { 'NOPROC' } else { 'CMD=' + $p.CommandLine; 'EXE=' + $p.ExecutablePath }`

	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script).Output()
	if err == nil {
		return parsePowerShellOutput(string(out))
	}

	out, wmicErr := exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline").Output()
	if wmicErr != nil {
		return processInfo{}, fmt.Errorf("powershell: %v; wmic: %v", err, wmicErr)
	}
	text := strings.ReplaceAll(string(out), "\r", "\n")
	if !strings.Contains(text, "--app-port") {
		return processInfo{}, ErrNotRunning
	}
	return processInfo{CommandLine: text}, nil
}

// parsePowerShellOutput interprets the output of the script in windowsProcess.
func parsePowerShellOutput(out string) (processInfo, error) {
	var info processInfo
	found := false
	for _, line := range strings.Split(strings.ReplaceAll(out, "\r", "\n"), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case line == "NOPROC":
			return processInfo{}, ErrNotRunning
		case strings.HasPrefix(line, "CMD="):
			found = true
			info.CommandLine = strings.TrimPrefix(line, "CMD=")
		case strings.HasPrefix(line, "EXE="):
			if exe := strings.TrimPrefix(line, "EXE="); exe != "" {
				info.InstallDir = filepath.Dir(exe)
			}
		}
	}
	if !found {
		return info, ErrNotRunning
	}
	if info.CommandLine == "" {
		return info, ErrAccessDenied
	}
	return info, nil
}
