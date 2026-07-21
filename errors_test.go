package clients3

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type stubAPIError struct {
	code string
}

func (e stubAPIError) Error() string            { return e.code }
func (e stubAPIError) ErrorCode() string        { return e.code }
func (e stubAPIError) ErrorMessage() string     { return e.code }
func (e stubAPIError) ErrorFault() smithy.ErrorFault {
	return smithy.FaultClient
}

func TestMapSDKNotFound(t *testing.T) {
	cases := []error{
		&types.NoSuchKey{},
		&types.NoSuchBucket{},
		&types.NotFound{},
		stubAPIError{code: "NoSuchKey"},
		stubAPIError{code: "NotFound"},
	}
	for _, in := range cases {
		err := mapSDK(in)
		if !IsNotFound(err) {
			t.Fatalf("expected IsNotFound for %T %v, got %v", in, in, err)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected wrap ErrNotFound: %v", err)
		}
	}
}

func TestMapSDKOther(t *testing.T) {
	in := fmt.Errorf("boom")
	err := mapSDK(in)
	if IsNotFound(err) {
		t.Fatal("unexpected not found")
	}
	if !errors.Is(err, in) && err != in {
		t.Fatalf("expected passthrough, got %v", err)
	}
}

func TestMapSDKNil(t *testing.T) {
	if mapSDK(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestCleanKey(t *testing.T) {
	k, err := cleanKey("/ui/auth/bg.png")
	if err != nil || k != "ui/auth/bg.png" {
		t.Fatalf("got %q %v", k, err)
	}
	if _, err := cleanKey(""); err == nil {
		t.Fatal("expected empty key error")
	}
	if _, err := cleanKey("../etc/passwd"); err == nil {
		t.Fatal("expected .. error")
	}
}
