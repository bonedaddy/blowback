//go:build !codeanalysis
// +build !codeanalysis

package blowback

import (
	"testing"

	"github.com/coredns/caddy"
)

// TestSetup tests the various things that should be parsed by setup.
// Make sure you also test for parse errors.
func TestSetup(t *testing.T) {
	c := caddy.NewTestController("dns", `blowback`)
	if err := setup(c); err == nil {
		t.Fatal("expected errors")
	}

	c = caddy.NewTestController("dns", `blowback proxy_server`)
	if err := setup(c); err == nil {
		t.Fatal("expected errors")
	}

	c = caddy.NewTestController("dns", `blowback proxy_server 0`)
	if err := setup(c); err == nil {
		t.Fatal("expected errors")
	}

	c = caddy.NewTestController("dns", `blowback proxy_server x`)
	if err := setup(c); err == nil {
		t.Fatal("expected errors")
	}

	c = caddy.NewTestController("dns", `blowback proxy_server http://127.0.0.1:8080`)
	if err := setup(c); err != nil {
		t.Fatalf("Expected no errors, but got: %v", err)
	}
}
