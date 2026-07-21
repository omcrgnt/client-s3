package clients3_test

import (
	"testing"

	clients3 "github.com/omcrgnt/client-s3"
)

func TestCredentialsStaticKeysRequired(t *testing.T) {
	if err := clients3.AccessKey("").Validate(); err == nil {
		t.Fatal("expected error for empty access key")
	}
	if err := clients3.SecretKey("").Validate(); err == nil {
		t.Fatal("expected error for empty secret key")
	}
}
