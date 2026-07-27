package clients3

import (
	"context"
	"fmt"
)

// ProbeReady reports S3 traffic readiness via HeadBucket (ops duck typing).
func (c *Client[C]) ProbeReady(ctx context.Context) error {
	if c == nil || c.s3 == nil {
		return fmt.Errorf("clients3: client not initialized")
	}
	return c.Ready(ctx)
}
