package utils

import (
	"context"
	"testing"
	"time"
)

func TestBuildCloudinaryUploadContext_UsesDedicatedCloudinaryTimeout(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	child, childCancel := buildCloudinaryUploadContext(parent)
	defer childCancel()

	deadline, ok := child.Deadline()
	if !ok {
		t.Fatal("expected child context to have a deadline")
	}

	remaining := time.Until(deadline)
	if remaining < 25*time.Second || remaining > 35*time.Second {
		t.Fatalf("expected upload timeout near 30s, got %v remaining", remaining)
	}
}
