package clients3

import (
	"context"
	"testing"
)

func TestProbeReadyUninitialized(t *testing.T) {
	if err := (*Client[*CredentialsStatic[Default]])(nil).ProbeReady(context.Background()); err == nil {
		t.Fatal("expected error for nil client")
	}
	c := &Client[*CredentialsStatic[Default]]{}
	if err := c.ProbeReady(context.Background()); err == nil {
		t.Fatal("expected error before Inject")
	}
}
