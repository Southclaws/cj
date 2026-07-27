package admin

import (
	"reflect"
	"testing"
)

func TestResolveAllowedHostsExplicitCSV(t *testing.T) {
	cfg := dashboardConfig{AllowedHosts: " Example.com , 100.64.1.2 ,,"}

	hosts, err := resolveAllowedHosts(cfg, "100.64.1.2:0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"example.com", "100.64.1.2"}
	if !reflect.DeepEqual(hosts, want) {
		t.Fatalf("got %v want %v", hosts, want)
	}
}

func TestResolveAllowedHostsDefaultsForImplicitLoopback(t *testing.T) {
	cfg := dashboardConfig{}

	hosts, err := resolveAllowedHosts(cfg, "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"localhost", "127.0.0.1"}
	if !reflect.DeepEqual(hosts, want) {
		t.Fatalf("got %v want %v", hosts, want)
	}
}

func TestResolveAllowedHostsDefaultsIncludeResolvedNonLoopbackAddress(t *testing.T) {
	cfg := dashboardConfig{}

	hosts, err := resolveAllowedHosts(cfg, "100.64.1.2:0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"localhost", "127.0.0.1", "100.64.1.2"}
	if !reflect.DeepEqual(hosts, want) {
		t.Fatalf("got %v want %v", hosts, want)
	}
}

func TestResolveAllowedHostsRequiresExplicitConfigWhenListenAddrIsExplicit(t *testing.T) {
	cfg := dashboardConfig{ListenAddr: "0.0.0.0:9999"}

	_, err := resolveAllowedHosts(cfg, "0.0.0.0:9999")
	if err == nil {
		t.Fatal("expected an error when listen addr is explicit but allowed hosts is not")
	}
}
