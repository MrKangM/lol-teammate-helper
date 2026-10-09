// Package lcu locates a running League client and extracts its API credentials.
package lcu

import (
	"errors"
	"fmt"
	"os/exec"
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

// Detect finds the running League client and returns its credentials.
func Detect() (Credentials, error) {
	cmdline, err := clientCommandLine()
	if err != nil {
		return Credentials{}, err
	}
	return ParseCommandLine(cmdline)
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

func clientCommandLine() (string, error) {
	if runtime.GOOS == "windows" {
		return windowsCommandLine()
	}
	out, err := exec.Command("ps", "-A", "-o", "args").Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "LeagueClientUx") {
			return line, nil
		}
	}
	return "", errors.New("league client is not running")
}

// windowsCommandLine tries PowerShell first and falls back to wmic, which is
// removed from recent Windows 11 builds.
func windowsCommandLine() (string, error) {
	const ps = "(Get-CimInstance Win32_Process -Filter \"name='LeagueClientUx.exe'\").CommandLine"
	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps).Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return string(out), nil
	}

	out, wmicErr := exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline").Output()
	if wmicErr != nil {
		if err == nil {
			err = errors.New("league client is not running")
		}
		return "", err
	}
	text := strings.ReplaceAll(string(out), "\r", "\n")
	if !strings.Contains(text, "--app-port") {
		return "", errors.New("league client is not running")
	}
	return text, nil
}
