package clients3_test

import (
	"context"
	"reflect"
	"testing"

	clients3 "github.com/omcrgnt/client-s3"
	httpv1 "github.com/omcrgnt/proto/gen/go/http/v1"
)

func TestClientBuildDepsInject(t *testing.T) {
	cfg := &clients3.Config[*clients3.CredentialsStatic[clients3.Default]]{
		Bucket:   "bemvpgame-assets",
		Endpoint: httpv1.URL{Value: "http://127.0.0.1:9000"},
		Region:   "us-east-1",
	}
	got, err := cfg.Build()
	if err != nil {
		t.Fatal(err)
	}
	c := got.(*clients3.Client[*clients3.CredentialsStatic[clients3.Default]])
	if c.Bucket() != "bemvpgame-assets" {
		t.Fatalf("bucket: %q", c.Bucket())
	}

	built, err := (&clients3.CredentialsStaticConfig[clients3.Default]{
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
	}).Build()
	if err != nil {
		t.Fatal(err)
	}
	creds := built.(*clients3.CredentialsStatic[clients3.Default])

	deps := c.Deps()
	if got, want := reflect.TypeOf(deps[0]), reflect.TypeOf((*clients3.CredentialsStatic[clients3.Default])(nil)); got != want {
		t.Fatalf("Deps()[0] type = %v, want %v", got, want)
	}

	c.Inject([]any{creds})

	cred, err := creds.Retrieve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if cred.AccessKeyID != "minioadmin" || cred.SecretAccessKey != "minioadmin" {
		t.Fatalf("creds: %+v", cred)
	}
}
