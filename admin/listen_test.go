package admin

import "testing"

func TestResolveListenAddrExplicitWins(t *testing.T) {
	cfg := dashboardConfig{ListenAddr: "0.0.0.0:9999"}

	addr, source := resolveListenAddr(cfg)

	if addr != "0.0.0.0:9999" || source != "explicit" {
		t.Fatalf("got addr=%q source=%q", addr, source)
	}
}

func TestResolveListenAddrDefaultsToLoopback(t *testing.T) {
	cfg := dashboardConfig{}

	addr, source := resolveListenAddr(cfg)

	if addr != "127.0.0.1:0" || source != "loopback" {
		t.Fatalf("got addr=%q source=%q", addr, source)
	}
}
