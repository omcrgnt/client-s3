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

// CredentialsStaticConfig is the ecfg spec for [CredentialsStatic].
type CredentialsStaticConfig struct {
	AccessKey AccessKey
	SecretKey SecretKey
}

// Build returns a [CredentialsStatic] resource.
func (c *CredentialsStaticConfig) Build() (any, error) {
	return &CredentialsStatic{
		accessKey: c.AccessKey.String(),
		secretKey: c.SecretKey.String(),
	}, nil
}

// CredentialsStatic is long-lived access/secret credentials (MinIO / IAM user keys).
// Catalog field: *CredentialsStatic (Configurable). Wire with Client[*CredentialsStatic].
type CredentialsStatic struct {
	accessKey string
	secretKey string
}

var (
	_ app.Configurable        = (*CredentialsStatic)(nil)
	_ CredentialsProvider     = (*CredentialsStatic)(nil)
	_ aws.CredentialsProvider = (*CredentialsStatic)(nil)
)

// BuildConfig returns the config spec for materialize.
func (*CredentialsStatic) BuildConfig() (app.Materializer, error) {
	return &CredentialsStaticConfig{}, nil
}

// Retrieve implements [aws.CredentialsProvider] / [CredentialsProvider].
func (c *CredentialsStatic) Retrieve(_ context.Context) (aws.Credentials, error) {
	if c == nil {
		return aws.Credentials{}, fmt.Errorf("clients3: credentials static is nil")
	}
	return aws.Credentials{
		AccessKeyID:     c.accessKey,
		SecretAccessKey: c.secretKey,
		Source:          "clients3.CredentialsStatic",
	}, nil
}
