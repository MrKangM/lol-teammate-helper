package config

import "testing"

func TestUpdateAndInstance(t *testing.T) {
	if !Update(1234, "tok", "HN1") {
		t.Fatal("first update should report a change")
	}
	if Update(1234, "tok", "HN1") {
		t.Fatal("identical update should not report a change")
	}
	if !Update(4321, "tok2", "HN1") {
		t.Fatal("new credentials should report a change")
	}

	cfg, ok := Instance()
	if !ok || cfg.Port != 4321 || cfg.MetaToken != "tok2" || cfg.Region != "艾欧尼亚" {
		t.Fatalf("unexpected config %+v ok=%v", cfg, ok)
	}
	cfg.Port = 1 // a snapshot must not alias the shared state
	if again, _ := Instance(); again.Port != 4321 {
		t.Fatal("Instance leaked internal state")
	}
}
