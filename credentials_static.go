package clients3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/omcrgnt/app"
)

// AccessKey is a static access key id.
type AccessKey string

func (AccessKey) Usage() string { return "Static access key id" }

func (a AccessKey) Validate() error {
	if strings.TrimSpace(string(a)) == "" {
		return fmt.Errorf("access key is required")
	}
	return nil
}

func (a AccessKey) String() string { return strings.TrimSpace(string(a)) }

// SecretKey is a static secret access key.
type SecretKey string

func (SecretKey) Usage() string { return "Static secret access key" }

func (s SecretKey) Validate() error {
	if strings.TrimSpace(string(s)) == "" {
		return fmt.Errorf("secret key is required")
	}
	return nil
}

func (s SecretKey) String() string { return strings.TrimSpace(string(s)) }

// Default tags the single static credentials set in a process (most apps).
// Also used when several clients share one credentials resource.
type Default struct{}

// CredentialsStaticConfig is the ecfg spec for [CredentialsStatic].
type CredentialsStaticConfig[Tag any] struct {
	AccessKey AccessKey
	SecretKey SecretKey
}

// Build returns a [CredentialsStatic] resource.
func (c *CredentialsStaticConfig[Tag]) Build() (any, error) {
	return &CredentialsStatic[Tag]{
		accessKey: c.AccessKey.String(),
		secretKey: c.SecretKey.String(),
	}, nil
}

// CredentialsStatic is long-lived access/secret credentials (MinIO / IAM user keys).
// Tag distinguishes multiple static sets in one process for SDI.
// One set: CredentialsStatic[Default]. Two sets: CredentialsStatic[assets], CredentialsStatic[backups], …
// Catalog: *CredentialsStatic[Tag]; wire Client[*CredentialsStatic[Tag]].
type CredentialsStatic[Tag any] struct {
	accessKey string
	secretKey string
}

var (
	_ app.Configurable        = (*CredentialsStatic[Default])(nil)
	_ CredentialsProvider     = (*CredentialsStatic[Default])(nil)
	_ aws.CredentialsProvider = (*CredentialsStatic[Default])(nil)
)

// BuildConfig returns the config spec for materialize.
func (*CredentialsStatic[Tag]) BuildConfig() (app.Materializer, error) {
	return &CredentialsStaticConfig[Tag]{}, nil
}

// Retrieve implements [aws.CredentialsProvider] / [CredentialsProvider].
func (c *CredentialsStatic[Tag]) Retrieve(_ context.Context) (aws.Credentials, error) {
	if c == nil {
		return aws.Credentials{}, fmt.Errorf("clients3: credentials static is nil")
	}
	return aws.Credentials{
		AccessKeyID:     c.accessKey,
		SecretAccessKey: c.secretKey,
		Source:          "clients3.CredentialsStatic",
	}, nil
}
