package clients3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/tracing/smithyoteltracing"
	"github.com/omcrgnt/app"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/otel"
)

// Client is an S3-compatible object client bound to one bucket and credentials type C.
// Catalog field: *Client[C] (Configurable); C is a concrete [CredentialsProvider]
// (e.g. *CredentialsStatic[Default]).
// SDI wires C via Deps/Inject (same pattern as srv-http.Server[T]).
type Client[C CredentialsProvider] struct {
	s3        *s3.Client
	presign   *s3.PresignClient
	bucket    string
	label     string
	endpoint  string
	region    string
	pathStyle bool
}

var _ app.Configurable = (*Client[*CredentialsStatic[Default]])(nil)

// BuildConfig returns the config spec for materialize.
func (*Client[C]) BuildConfig() (app.Materializer, error) {
	return &Config[C]{}, nil
}

// Build returns a Client resource; AWS SDK is wired in Inject after credentials resolve.
func (cfg *Config[C]) Build() (any, error) {
	return &Client[C]{
		bucket:    cfg.Bucket.String(),
		label:     cfg.Label.GetValue(),
		endpoint:  cfg.Endpoint.GetValue(),
		region:    cfg.Region.String(),
		pathStyle: true, // required for MinIO / custom endpoint
	}, nil
}

// Deps declares the concrete credentials implementor C.
func (c *Client[C]) Deps() []any {
	var zero C
	return []any{zero}
}

// Inject receives credentials and materializes the AWS SDK client.
func (c *Client[C]) Inject(args []any) {
	var cred C
	found := false
	for _, arg := range args {
		if v, ok := arg.(C); ok {
			cred = v
			found = true
			break
		}
	}
	if !found {
		return
	}

	awsCfg := aws.Config{
		Region:      c.region,
		Credentials: cred,
	}
	otelaws.AppendMiddlewares(&awsCfg.APIOptions)

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.TracerProvider = smithyoteltracing.Adapt(otel.GetTracerProvider())
		o.BaseEndpoint = aws.String(c.endpoint)
		o.UsePathStyle = c.pathStyle
	})
	c.s3 = client
	c.presign = s3.NewPresignClient(client)
}

// Label returns the configured resource label.
func (c *Client[C]) Label() string {
	if c == nil {
		return ""
	}
	return c.label
}

// Bucket returns the bucket this client is bound to.
func (c *Client[C]) Bucket() string {
	if c == nil {
		return ""
	}
	return c.bucket
}

// NewTestClient constructs a Client around an existing SDK client (tests only).
func NewTestClient[C CredentialsProvider](inner *s3.Client, bucket string) *Client[C] {
	return &Client[C]{
		s3:      inner,
		presign: s3.NewPresignClient(inner),
		bucket:  bucket,
	}
}

// Close is a no-op; retained for lifecycle symmetry with other resources.
func (c *Client[C]) Close(_ context.Context) error {
	return nil
}

// PutOptions controls optional PutObject fields.
type PutOptions struct {
	ContentType   string
	ContentLength int64 // 0 = omit (SDK may buffer)
}

// Head holds object metadata from HeadObject / GetObject.
type Head struct {
	ContentType   string
	ContentLength int64
	ETag          string
	SHA256        string
}

// Put uploads an object at key.
func (c *Client[C]) Put(ctx context.Context, key string, body io.Reader, opts PutOptions) error {
	key, err := cleanKey(key)
	if err != nil {
		return err
	}

	in := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if opts.ContentType != "" {
		in.ContentType = aws.String(opts.ContentType)
	}
	if opts.ContentLength > 0 {
		in.ContentLength = aws.Int64(opts.ContentLength)
	}

	_, err = c.s3.PutObject(ctx, in)
	return mapSDK(err)
}

// Get downloads an object. Caller must Close the returned body.
func (c *Client[C]) Get(ctx context.Context, key string) (io.ReadCloser, Head, error) {
	key, err := cleanKey(key)
	if err != nil {
		return nil, Head{}, err
	}

	out, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, Head{}, mapSDK(err)
	}
	return out.Body, headFromGet(out), nil
}

// Head fetches object metadata without the body.
func (c *Client[C]) Head(ctx context.Context, key string) (Head, error) {
	key, err := cleanKey(key)
	if err != nil {
		return Head{}, err
	}

	out, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return Head{}, mapSDK(err)
	}
	return headFromHead(out), nil
}

// Delete removes an object.
func (c *Client[C]) Delete(ctx context.Context, key string) error {
	key, err := cleanKey(key)
	if err != nil {
		return err
	}

	_, err = c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return mapSDK(err)
}

// PresignGet returns a time-limited GET URL for key.
func (c *Client[C]) PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error) {
	key, err := cleanKey(key)
	if err != nil {
		return "", err
	}
	if ttl <= 0 {
		return "", fmt.Errorf("clients3: ttl must be positive")
	}

	out, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", mapSDK(err)
	}
	return out.URL, nil
}

// Ready checks that the configured bucket exists and is reachable.
func (c *Client[C]) Ready(ctx context.Context) error {
	_, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.bucket),
	})
	return mapSDK(err)
}

func cleanKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return "", fmt.Errorf("clients3: key is required")
	}
	if strings.Contains(key, "..") {
		return "", fmt.Errorf("clients3: key must not contain ..")
	}
	return key, nil
}

func headFromGet(out *s3.GetObjectOutput) Head {
	h := Head{}
	if out.ContentType != nil {
		h.ContentType = *out.ContentType
	}
	if out.ContentLength != nil {
		h.ContentLength = *out.ContentLength
	}
	if out.ETag != nil {
		h.ETag = strings.Trim(*out.ETag, `"`)
	}
	h.SHA256 = checksumSHA256(out.ChecksumSHA256, out.Metadata)
	return h
}

func headFromHead(out *s3.HeadObjectOutput) Head {
	h := Head{}
	if out.ContentType != nil {
		h.ContentType = *out.ContentType
	}
	if out.ContentLength != nil {
		h.ContentLength = *out.ContentLength
	}
	if out.ETag != nil {
		h.ETag = strings.Trim(*out.ETag, `"`)
	}
	h.SHA256 = checksumSHA256(out.ChecksumSHA256, out.Metadata)
	return h
}

func checksumSHA256(sum *string, meta map[string]string) string {
	if sum != nil && *sum != "" {
		return *sum
	}
	if meta != nil {
		if v := meta["sha256"]; v != "" {
			return v
		}
		if v := meta["SHA256"]; v != "" {
			return v
		}
	}
	return ""
}
