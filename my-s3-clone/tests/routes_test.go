package tests

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gorilla/mux"
	"my-s3-clone/handlers"
	"my-s3-clone/router"
	"my-s3-clone/dto"
	"io"
	"time"
	"fmt"
)

// Mock implementation of FileInfo 
type MockFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return m.size }
func (m MockFileInfo) Mode() os.FileMode  { return 0644 }
func (m MockFileInfo) ModTime() time.Time { return m.modTime }
func (m MockFileInfo) IsDir() bool        { return false }
func (m MockFileInfo) Sys() interface{}   { return nil }

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	AddObjectFunc         func(bucketName, objectName string, data io.Reader, contentSha256 string) error
	DeleteObjectFunc      func(bucketName, objectName string) error
	CheckBucketExistsFunc func(bucketName string) (bool, error)
	CheckObjectExistFunc  func(bucketName, objectName string) (bool, time.Time, int64, error)
	DeleteBucketFunc      func(bucketName string) error
	GetObjectFunc         func(bucketName, objectName string) ([]byte, dto.FileInfo, error)
	ListBucketsFunc       func() []string
	ListObjectsFunc       func(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error)
	CreateBucketFunc      func(bucketName string) error
}

// Implementations of the Storage interface using the mock functions
func (m *MockStorage) AddObject(bucketName, objectName string, data io.Reader, contentSha256 string) error {
	if m.AddObjectFunc != nil {
		return m.AddObjectFunc(bucketName, objectName, data, contentSha256)
	}
	return nil
}

func (m *MockStorage) DeleteObject(bucketName, objectName string) error {
	if m.DeleteObjectFunc != nil {
		return m.DeleteObjectFunc(bucketName, objectName)
	}
	return nil
}

func (m *MockStorage) CheckBucketExists(bucketName string) (bool, error) {
	if m.CheckBucketExistsFunc != nil {
		return m.CheckBucketExistsFunc(bucketName)
	}
	return false, nil
}

func (m *MockStorage) CheckObjectExist(bucketName, objectName string) (bool, time.Time, int64, error) {
	if m.CheckObjectExistFunc != nil {
		return m.CheckObjectExistFunc(bucketName, objectName)
	}
	return false, time.Time{}, 0, nil
}

func (m *MockStorage) DeleteBucket(bucketName string) error {
	if m.DeleteBucketFunc != nil {
		return m.DeleteBucketFunc(bucketName)
	}
	return nil
}

func (m *MockStorage) GetObject(bucketName, objectName string) ([]byte, dto.FileInfo, error) {
	if m.GetObjectFunc != nil {
		return m.GetObjectFunc(bucketName, objectName)
	}
	return nil, nil, os.ErrNotExist
}

func (m *MockStorage) ListBuckets() []string {
	if m.ListBucketsFunc != nil {
		return m.ListBucketsFunc()
	}
	return []string{}
}

func (m *MockStorage) ListObjects(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error) {
	if m.ListObjectsFunc != nil {
		return m.ListObjectsFunc(bucketName, prefix, marker, maxKeys)
	}
	return dto.ListObjectsResponse{}, nil
}

// Mock implementation of CreateBucket
func (m *MockStorage) CreateBucket(bucketName string) error {
    if m.CreateBucketFunc != nil {
        return m.CreateBucketFunc(bucketName)
    }
    return nil
}

// Test for the /probe-bsign{suffix:.*} route
func TestProbeBSignRoute(t *testing.T) {
	r := router.SetupRouter()

	tests := []struct {
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{"GET", "/probe-bsign", http.StatusOK, "<Response></Response>"},
		{"HEAD", "/probe-bsign", http.StatusOK, ""},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d", tt.expectedCode, rr.Code)
		}

		if tt.method == "GET" && rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q but got %q", tt.expectedBody, rr.Body.String())
		}
	}
}

// Test for the /{bucketName}/?delete= (POST batch delete)
func TestHandleDeleteObject(t *testing.T) {
	// Mock storage
	mockStorage := &MockStorage{}

	// Initialize the router with mock storage
	r := router.SetupRouterWithStorage(mockStorage) // Assuming you have a way to inject storage into the router

	// Create a DeleteBatchRequest with objects to delete
	deleteReq := dto.DeleteObjectRequest{
		Objects: []dto.ObjectToDelete{
			{Key: "object1.txt"},
			{Key: "object2.txt"},
		},
	}

	// Convert the request to XML
	body, err := xml.Marshal(deleteReq)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}

	// Create the POST request with the XML body
	req, err := http.NewRequest("POST", "/bucketName/?delete=", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the router with the simulated request
	r.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response content
	var deleteResult dto.DeleteResult
	err = xml.Unmarshal(rr.Body.Bytes(), &deleteResult)
	if err != nil {
		t.Fatalf("Error unmarshaling response body: %v", err)
	}

	// Check that the deleted objects match those in the request
	expectedDeleted := []string{"object1.txt", "object2.txt"}
	if len(deleteResult.DeletedResult) != len(expectedDeleted) {
		t.Errorf("expected %d deleted objects, got %d", len(expectedDeleted), len(deleteResult.DeletedResult))
	}

	for i, obj := range deleteResult.DeletedResult {
		if obj.Key != expectedDeleted[i] {
			t.Errorf("expected deleted object %s, got %s", expectedDeleted[i], obj.Key)
		}
	}
}

// Test for the /{bucketName}/{objectName} (POST/PUT) route
func TestHandleAddObject(t *testing.T) {
	// Create a new instance of the mock storage
	mockStorage := &MockStorage{
		AddObjectFunc: func(bucketName, objectName string, data io.Reader, contentSha256 string) error {
			if bucketName == "test-bucket" && objectName == "test-object" {
				// Simulate successful upload, reading the content from the reader
				buf := new(bytes.Buffer)
				if _, err := buf.ReadFrom(data); err != nil {
					return err
				}
				if buf.String() != "file content" {
					return fmt.Errorf("unexpected file content: %s", buf.String())
				}
				return nil
			}
			return os.ErrNotExist // Simulate failure
		},
		CheckBucketExistsFunc: func(bucketName string) (bool, error) {
			if bucketName == "test-bucket" {
				return true, nil
			}
			return false, nil
		},
		CheckObjectExistFunc: func(bucketName, objectName string) (bool, time.Time, int64, error) {
			if bucketName == "test-bucket" && objectName == "test-object" {
				return true, time.Now(), 1234, nil
			}
			return false, time.Time{}, 0, os.ErrNotExist
		},
	}

	// Initialize the router with the mock storage
	r := mux.NewRouter()
	r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleAddObject(mockStorage)).Methods("POST", "PUT")

	// Create a POST request to upload an object
	req, err := http.NewRequest("POST", "/test-bucket/test-object", bytes.NewBuffer([]byte("file content")))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("X-Amz-Content-Sha256", "dummyhash")
	req.Header.Set("Expect", "100-continue") // Add the Expect header
	req.Header.Set("X-Amz-Decoded-Content-Length", "12") // Add the missing header

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check for the 100 Continue response if Expect: 100-continue is set
	if rr.Result().StatusCode == http.StatusContinue {
		rr = httptest.NewRecorder() // Reset response recorder for actual processing
		r.ServeHTTP(rr, req)        // Process the actual request after 100 Continue
	}

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := ""
	if rr.Body.String() != expectedResponse {
		t.Errorf("expected body %q but got %q", expectedResponse, rr.Body.String())
	}

	// Validate the response headers
	if rr.Header().Get("ETag") == "" {
		t.Errorf("expected ETag header to be set")
	}
	if rr.Header().Get("x-amz-id-2") == "" {
		t.Errorf("expected x-amz-id-2 header to be set")
	}
	if rr.Header().Get("x-amz-request-id") == "" {
		t.Errorf("expected x-amz-request-id header to be set")
	}
}

func TestHandleCheckObjectExist(t *testing.T) {
	// Create a new instance of the mock storage
	mockStorage := &MockStorage{
		CheckObjectExistFunc: func(bucketName, objectName string) (bool, time.Time, int64, error) {
			// Simulate that the object exists
			if bucketName == "test-bucket" && objectName == "test-object" {
				return true, time.Now(), 1234, nil
			}
			// Simulate that the object does not exist
			return false, time.Time{}, 0, nil
		},
	}

	// Initialize the router with the mock storage
	r := mux.NewRouter()
	r.HandleFunc("/{bucketName}/{objectName}", handlers.HandleCheckObjectExist(mockStorage)).Methods("HEAD")

	// Create test cases
	tests := []struct {
		bucketName   string
		objectName   string
		expectedCode int
	}{
		{"test-bucket", "test-object", http.StatusOK},       
		{"test-bucket", "nonexistent-object", http.StatusNotFound}, 
	}

	for _, tt := range tests {
		// Create a new HEAD request for the object
		req, err := http.NewRequest("HEAD", "/"+tt.bucketName+"/"+tt.objectName, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Serve the request using the router
		r.ServeHTTP(rr, req)

		// Check the status code
		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d for object: %s", tt.expectedCode, rr.Code, tt.objectName)
		}

		// Additional checks if the object exists
		if tt.expectedCode == http.StatusOK {
			// Check for the Last-Modified and Content-Length headers
			if rr.Header().Get("Last-Modified") == "" {
				t.Errorf("expected Last-Modified header to be set for object: %s", tt.objectName)
			}
			if rr.Header().Get("Content-Length") != "1234" {
				t.Errorf("expected Content-Length to be 1234 but got %s", rr.Header().Get("Content-Length"))
			}
		}
	}
}

func TestHandleListObjects(t *testing.T) {
    // Create a new instance of the mock storage
    mockStorage := &MockStorage{
        ListObjectsFunc: func(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error) {
            // Simulate a response with some objects
            if marker == "object1.txt" {
                // Simulate paginated response
                return dto.ListObjectsResponse{
                    Name:    bucketName,
                    Prefix:  prefix,
                    Marker:  marker,
                    MaxKeys: maxKeys,
                    Contents: []dto.Object{
                        {Key: "object2.txt", LastModified: time.Now(), Size: 5678},
                    },
                    IsTruncated: false,
                }, nil
            }

            return dto.ListObjectsResponse{
                Name:    bucketName,
                Prefix:  prefix,
                Marker:  marker,
                MaxKeys: maxKeys,
                Contents: []dto.Object{
                    {Key: "object1.txt", LastModified: time.Now(), Size: 1234},
                    {Key: "object2.txt", LastModified: time.Now(), Size: 5678},
                },
                IsTruncated: false,
            }, nil
        },
    }

    // Initialize the router with the mock storage
    r := mux.NewRouter()
    r.HandleFunc("/{bucketName}/", handlers.HandleListObjects(mockStorage)).Methods("GET", "HEAD")

    // Create test cases
    tests := []struct {
        bucketName   string
        queryParams  string
        expectedCode int
        expectedKeys []string
    }{
        {
            "test-bucket", "prefix=&marker=&max-keys=2", http.StatusOK,
            []string{"object1.txt", "object2.txt"},
        },
        {
            "test-bucket", "prefix=&marker=object1.txt&max-keys=1", http.StatusOK,
            []string{"object2.txt"},
        },
    }

    for _, tt := range tests {
        // Create a GET request for listing objects in the bucket
        req, err := http.NewRequest("GET", "/"+tt.bucketName+"/?"+tt.queryParams, nil)
        if err != nil {
            t.Fatalf("could not create request: %v", err)
        }

        // Create a response recorder to capture the response
        rr := httptest.NewRecorder()

        // Serve the request using the router
        r.ServeHTTP(rr, req)

        // Check the status code
        if rr.Code != tt.expectedCode {
            t.Errorf("expected status %d but got %d for bucket: %s", tt.expectedCode, rr.Code, tt.bucketName)
        }

        // Check the response content if the request was successful
        if tt.expectedCode == http.StatusOK {
            var listObjectsResponse dto.ListObjectsResponse
            err := xml.Unmarshal(rr.Body.Bytes(), &listObjectsResponse)
            if err != nil {
                t.Fatalf("Error unmarshaling response body: %v", err)
            }

            // Check the number of returned objects
            if len(listObjectsResponse.Contents) != len(tt.expectedKeys) {
                t.Errorf("expected %d objects but got %d", len(tt.expectedKeys), len(listObjectsResponse.Contents))
            }

            // Check that the object keys match the expected keys
            for i, obj := range listObjectsResponse.Contents {
                if obj.Key != tt.expectedKeys[i] {
                    t.Errorf("expected object key %s but got %s", tt.expectedKeys[i], obj.Key)
                }
            }
        }
    }
}

// Test for HandleGetBucket
func TestHandleGetBucket(t *testing.T) {
	// Mock storage
	mockStorage := &MockStorage{
		CheckBucketExistsFunc: func(bucketName string) (bool, error) {
			if bucketName == "existing-bucket" {
				return true, nil
			}
			return false, nil
		},
	}

	// Initialize the router with mock storage
	r := mux.NewRouter()
	r.HandleFunc("/{bucketName}/", handlers.HandleGetBucket(mockStorage)).Methods("GET")

	tests := []struct {
		bucketName   string
		location     string
		expectedCode int
		expectedBody string
	}{
		{"existing-bucket", "", http.StatusOK, "Bucket 'existing-bucket' exists and is accessible."},
		{"existing-bucket", "location", http.StatusOK, "<LocationConstraint>us-east-1</LocationConstraint>"},
		{"nonexistent-bucket", "", http.StatusNotFound, "Bucket 'nonexistent-bucket' not found\n"},
	}

	for _, tt := range tests {
		// Create a GET request
		reqURL := "/" + tt.bucketName + "/"
		if tt.location != "" {
			reqURL += "?location=" + tt.location
		}
		req, err := http.NewRequest("GET", reqURL, nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Serve the request using the router
		r.ServeHTTP(rr, req)

		// Check the status code
		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d for bucket: %s", tt.expectedCode, rr.Code, tt.bucketName)
		}

		// Check the response body
		if rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q but got %q", tt.expectedBody, rr.Body.String())
		}
	}
}

func TestHandleCreateBucket(t *testing.T) {
	// Mock storage
	mockStorage := &MockStorage{
		CreateBucketFunc: func(bucketName string) error {
			if bucketName == "test-bucket" {
				// Simulate successful bucket creation
				return nil
			}
			// Simulate an error if the bucket creation fails
			return fmt.Errorf("failed to create bucket")
		},
	}

	// Initialize the router with the mock storage
	r := mux.NewRouter()
	r.HandleFunc("/{bucketName}/", handlers.HandleCreateBucket(mockStorage)).Methods("PUT")

	tests := []struct {
		bucketName   string
		expectedCode int
		expectedBody string
	}{
		{"test-bucket", http.StatusOK, ""},         
		{"fail-bucket", http.StatusInternalServerError, "failed to create bucket\n"}, 
	}

	for _, tt := range tests {
		// Create a PUT request to create a bucket
		req, err := http.NewRequest("PUT", "/"+tt.bucketName+"/", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Serve the request using the router
		r.ServeHTTP(rr, req)

		// Check the status code
		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d for bucket: %s", tt.expectedCode, rr.Code, tt.bucketName)
		}

		// Check the response body

		if tt.expectedCode == http.StatusOK {
			// Définir la réponse XML attendue
			xmlResponse := `<ListAllMyBucketsResult><Buckets><Bucket><Name>test-bucket</Name><CreationDate>0001-01-01T00:00:00Z</CreationDate></Bucket></Buckets></ListAllMyBucketsResult>`
	
			// Récupérer la réponse brute depuis le body de la requête
			actualResponse := rr.Body.String()  // Convertir les bytes en string pour la comparaison
	
			// Comparer directement la réponse reçue avec la réponse attendue
			if actualResponse != xmlResponse {
				t.Errorf("Expected XML response to be: %s, but got: %s", xmlResponse, actualResponse)
			}
		} else if rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q but got %q", tt.expectedBody, rr.Body.String())
		}
	}
}

// TestHandleDeleteBucket tests the deletion of a bucket
func TestHandleDeleteBucket(t *testing.T) {
	// Set up the mock storage with DeleteBucket behavior
	mockStorage := &MockStorage{
		DeleteBucketFunc: func(bucketName string) error {
			// Simulate failure for a specific bucket
			if bucketName == "fail-bucket" {
				return fmt.Errorf("Failed to delete bucket\n")
			}
			return nil
		},
	}

	// Initialize the router with the handler and the mock storage
	r := mux.NewRouter()
	r.HandleFunc("/{bucketName}/", handlers.HandleDeleteBucket(mockStorage)).Methods("DELETE")

	// Test cases
	tests := []struct {
		bucketName    string
		expectedCode  int
		expectedBody  string
	}{
		{
			bucketName:   "test-bucket",
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			bucketName:   "fail-bucket",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to delete bucket\n",
		},
	}

	// Run through test cases
	for _, tt := range tests {
		// Create a DELETE request for the bucket
		req, err := http.NewRequest("DELETE", "/"+tt.bucketName+"/", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// Serve the request using the router
		r.ServeHTTP(rr, req)

		// Check the status code
		if rr.Code != tt.expectedCode {
			t.Errorf("expected status %d but got %d for bucket: %s", tt.expectedCode, rr.Code, tt.bucketName)
		}

		// Check the response body
		if rr.Body.String() != tt.expectedBody {
			t.Errorf("expected body %q but got %q for bucket: %s", tt.expectedBody, rr.Body.String(), tt.bucketName)
		}
	}
}

func TestHandleListBuckets(t *testing.T) {
	// Set up the mock storage
	mockStorage := &MockStorage{
		ListBucketsFunc: func() []string {
			// Return some bucket names
			return []string{"bucket1", "bucket2", "bucket3"}
		},
	}

	// Initialize the router with the handler and the mock storage
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HandleListBuckets(mockStorage)).Methods("GET")

	// Create a GET request to list buckets
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	r.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, rr.Code)
	}

	// Parse the response body
	var listBucketsResult dto.ListAllMyBucketsResult
	err = xml.NewDecoder(bytes.NewBuffer(rr.Body.Bytes())).Decode(&listBucketsResult)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	// Check that the response contains the correct buckets
	expectedBuckets := []string{"bucket1", "bucket2", "bucket3"}
	if len(listBucketsResult.Buckets) != len(expectedBuckets) {
		t.Fatalf("expected %d buckets but got %d", len(expectedBuckets), len(listBucketsResult.Buckets))
	}

	for i, bucket := range listBucketsResult.Buckets {
		if bucket.Name != expectedBuckets[i] {
			t.Errorf("expected bucket name %q but got %q", expectedBuckets[i], bucket.Name)
		}
	}

	// Check that the response content type is set to "application/xml"
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/xml" {
		t.Errorf("expected content type application/xml but got %q", contentType)
	}
}