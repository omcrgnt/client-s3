package clients3_test

import (
	"testing"

	"buf.build/go/protovalidate"
	clients3 "github.com/omcrgnt/client-s3"
	httpv1 "github.com/omcrgnt/proto/gen/go/http/v1"
)

func TestConfigBucketRequired(t *testing.T) {
	if err := clients3.Bucket("").Validate(); err == nil {
		t.Fatal("expected error for empty bucket")
	}
}

func TestHTTPURLRejectsInvalid(t *testing.T) {
	v, err := protovalidate.New()
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"", "not-a-url", "ws://127.0.0.1/ws"} {
		if err := v.Validate(&httpv1.URL{Value: value}); err == nil {
			t.Fatalf("expected protovalidate error for %q", value)
		}
	}
}
