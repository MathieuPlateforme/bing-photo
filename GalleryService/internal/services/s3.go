package services

import (
	"bytes"
	"fmt"
	"net/http"
)

// S3Service gère la communication avec l'API S3-like
type S3Service struct {
	APIURL string
}

// NewS3Service initialise un S3Service
func NewS3Service(apiURL string) *S3Service {
	return &S3Service{APIURL: apiURL}
}

// CreateFolder crée un dossier dans l'API S3-like
func (s *S3Service) CreateBucket(folderPath string) error {
    // Ajouter un "/" pour simuler un dossier
    url := fmt.Sprintf("%s/%s/", s.APIURL, folderPath)

    req, err := http.NewRequest("PUT", url, bytes.NewReader([]byte{}))
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/octet-stream") // Type MIME pour un objet vide

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create folder, status: %s", resp.Status)
    }

    return nil
}

