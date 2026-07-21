package clients3

import (
	"fmt"
	"strings"

	common "github.com/omcrgnt/proto/gen/go/common/v1"
	httpv1 "github.com/omcrgnt/proto/gen/go/http/v1"
)

// Region is the AWS / S3 region name.
type Region string

func (Region) Usage() string { return "S3 region (default us-east-1)" }

func (r Region) Validate() error { return nil }

func (r Region) String() string {
	s := strings.TrimSpace(string(r))
	if s == "" {
		return "us-east-1"
	}
	return s
}

// Bucket is the object bucket bound to this client.
type Bucket string

func (Bucket) Usage() string { return "S3 bucket name" }

func (b Bucket) Validate() error {
	if strings.TrimSpace(string(b)) == "" {
		return fmt.Errorf("bucket is required")
	}
	return nil
}

func (b Bucket) String() string {
	return strings.TrimSpace(string(b))
}

// UsePathStyle enables path-style addressing (required for typical MinIO setups).
type UsePathStyle bool

func (UsePathStyle) Usage() string {
	return "Use path-style S3 addressing (true for MinIO)"
}

func (UsePathStyle) Validate() error { return nil }

// Config is the client-s3 spec; ecfg fills before Build.
// Endpoint is protovalidated by ecfg as http.v1.URL (non-empty http/https).
// Credentials come from SDI via Client[C] → CredentialsProvider implementor C.
type Config[C CredentialsProvider] struct {
	Label        common.Label
	Endpoint     httpv1.URL
	Region       Region
	Bucket       Bucket
	UsePathStyle UsePathStyle
}
