package clients3

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// ErrNotFound is returned when the object (or bucket for Ready) does not exist.
var ErrNotFound = errors.New("clients3: not found")

// IsNotFound reports whether err is or wraps [ErrNotFound].
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func mapSDK(err error) error {
	if err == nil {
		return nil
	}
	if isNotFoundSDK(err) {
		return fmt.Errorf("%w: %v", ErrNotFound, err)
	}
	return err
}

func isNotFoundSDK(err error) bool {
	if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
		return true
	}
	if _, ok := errors.AsType[*types.NoSuchBucket](err); ok {
		return true
	}
	if _, ok := errors.AsType[*types.NotFound](err); ok {
		return true
	}
	if apiErr, ok := errors.AsType[smithy.APIError](err); ok {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchKey", "NoSuchBucket", "404":
			return true
		}
	}
	return false
}
