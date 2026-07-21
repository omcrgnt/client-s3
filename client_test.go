package clients3

import (
	"context"
	"testing"
)

func TestPresignRequiresPositiveTTL(t *testing.T) {
	c := &Client[*CredentialsStatic]{bucket: "b"}
	_, err := c.PresignGet(context.Background(), "k", 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCleanKeyRequired(t *testing.T) {
	c := &Client[*CredentialsStatic]{bucket: "b"}
	err := c.Put(context.Background(), "", nil, PutOptions{})
	if err == nil {
		t.Fatal("expected key error")
	}
}
