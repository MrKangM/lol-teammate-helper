package lcu

import "testing"

func TestParseCommandLine(t *testing.T) {
	line := `"C:\Riot Games\League of Legends\LeagueClientUx.exe" "--riotclient-auth-token=abc" "--remoting-auth-token=tok_en-12" "--app-port=54321" "--rso_platform_id=HN1" "--region=TENCENT"`
	got, err := ParseCommandLine(line)
	if err != nil {
		t.Fatal(err)
	}
	want := Credentials{Port: 54321, Token: "tok_en-12", Region: "HN1"}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestParseCommandLineErrors(t *testing.T) {
	cases := map[string]string{
		"no port":  `--remoting-auth-token=x`,
		"no token": `--app-port=1234`,
	}
	for name, line := range cases {
		if _, err := ParseCommandLine(line); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestParseCommandLineRegionOptional(t *testing.T) {
	got, err := ParseCommandLine(`--app-port=1 --remoting-auth-token=t`)
	if err != nil || got.Region != "" {
		t.Fatalf("got %+v, err %v", got, err)
	}
}

func TestParseLockfile(t *testing.T) {
	got, err := ParseLockfile("LeagueClient:1234:54321:secret-pw:https\n")
	if err != nil || got.Port != 54321 || got.Token != "secret-pw" {
		t.Fatalf("got %+v, err %v", got, err)
	}
	for _, bad := range []string{"", "a:b:c", "LeagueClient:1:notaport:pw:https", "LeagueClient:1:2::https"} {
		if _, err := ParseLockfile(bad); err == nil {
			t.Errorf("%q should fail", bad)
		}
	}
}

func TestParsePowerShellOutput(t *testing.T) {
	info, err := parsePowerShellOutput("CMD=\"x\" --app-port=1\r\nEXE=C:\\Riot Games\\LoL\\LeagueClientUx.exe\r\n")
	if err != nil || info.CommandLine == "" || info.InstallDir == "" {
		t.Fatalf("got %+v, err %v", info, err)
	}
	if _, err := parsePowerShellOutput("NOPROC\r\n"); err != ErrNotRunning {
		t.Errorf("NOPROC: got %v", err)
	}
	// process visible but command line hidden: the elevated-client case
	if _, err := parsePowerShellOutput("CMD=\r\nEXE=\r\n"); err != ErrAccessDenied {
		t.Errorf("empty command line: got %v", err)
	}
	if _, err := parsePowerShellOutput(""); err != ErrNotRunning {
		t.Errorf("empty output: got %v", err)
	}
}
