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
