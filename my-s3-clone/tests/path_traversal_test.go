package tests

import (
	"my-s3-clone/storage"
	"os"
    "path/filepath"
	"strings"
	"testing"
)

func TestPathTraversal(t *testing.T) {
	// Setup temporary storage root
	tmpDir, err := os.MkdirTemp("", "s3-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

    // Create a new FileStorage with the temp dir as root
    fs := storage.NewFileStorage(tmpDir)

    // Test case 1: Try to create a file outside the root
    bucketName := "bucket"
    objectName := "../outside_file"

    // We expect this to fail
    err = fs.AddObject(bucketName, objectName, strings.NewReader("content"), "")
    if err == nil {
        t.Errorf("Expected error for path traversal attempt, got nil")
    } else {
        if !strings.Contains(err.Error(), "path traversal attempt") {
             t.Errorf("Expected 'path traversal attempt' error, got: %v", err)
        }
    }

    // Test case 2: Try to read a file outside the root
    _, _, err = fs.GetObject(bucketName, "../../../etc/passwd")
    if err == nil {
        t.Errorf("Expected error for path traversal attempt, got nil")
    } else {
        if !strings.Contains(err.Error(), "path traversal attempt") {
             t.Errorf("Expected 'path traversal attempt' error, got: %v", err)
        }
    }

    // Test case 3: Partial Path Traversal
    // If root is /tmp/s3-test, try to access /tmp/s3-test-suffix
    rootBase := filepath.Base(tmpDir)
    siblingDir := rootBase + "-suffix"

    err = fs.AddObject("..", siblingDir+"/file", strings.NewReader("content"), "")
    if err == nil {
         t.Errorf("Expected error for partial path traversal attempt, got nil")
    } else {
         if !strings.Contains(err.Error(), "path traversal attempt") {
             t.Errorf("Expected 'path traversal attempt' error, got: %v", err)
         }
    }


    // Test case 4: Valid creation
    err = fs.CreateBucket(bucketName)
    if err != nil {
        t.Fatalf("Failed to create bucket: %v", err)
    }

    objectName = "valid_file.txt"
    err = fs.AddObject(bucketName, objectName, strings.NewReader("content"), "")
    if err != nil {
        t.Errorf("Failed to add valid object: %v", err)
    }
}
