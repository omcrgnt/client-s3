package clients3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// CredentialsProvider is the SDI port for S3 credentials (static, Vault, STS, …).
// Catalog: declare a concrete implementor (e.g. *CredentialsStatic) and bind Client[C] to that type.
type CredentialsProvider interface {
	Retrieve(ctx context.Context) (aws.Credentials, error)
}
