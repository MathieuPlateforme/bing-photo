package services

import (
	"bytes"
	"fmt"
	"net/http"
    "encoding/xml"
    "time"
	"io"
)

// S3Service gère la communication avec l'API S3-like
type S3Service struct {
	APIURL string
}

type ListAllMyBucketsResult struct {
    XMLName xml.Name `xml:"ListAllMyBucketsResult"`
    Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
    Name         string    `xml:"Name"`
    CreationDate time.Time `xml:"CreationDate"`
    LocationConstraint   string   `xml:"LocationConstraint,omitempty"`
    ObjectLockConfig   string   `xml:"ObjectLockConfiguration,omitempty"`
    ObjectDelimiter   string   `xml:"ObjectDelimiter,omitempty"`
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

// ListBuckets récupère la liste des buckets depuis l'API S3-like
func (s *S3Service) ListBuckets() ([]Bucket, error) {
	// Construire l'URL pour lister les buckets
	url := fmt.Sprintf("%s/", s.APIURL)

	// Envoyer une requête HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'appel à l'API S3-like : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de réponse
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("échec de la récupération des buckets, statut : %s", resp.Status)
	}

	// Décoder la réponse XML
	var bucketsResponse ListAllMyBucketsResult
	err = xml.NewDecoder(resp.Body).Decode(&bucketsResponse)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de la réponse XML : %v", err)
	}

	return bucketsResponse.Buckets, nil
}

func (s *S3Service) DeleteBucket(bucketName string) error {
	// Construire l'URL pour supprimer le bucket
	url := fmt.Sprintf("%s/%s/", s.APIURL, bucketName)

	// Envoyer une requête HTTP DELETE
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("échec de la création de la requête de suppression : %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("échec de la suppression du bucket : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de réponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("échec de la suppression du bucket, statut : %s", resp.Status)
	}

	return nil
}

func (s *S3Service) UploadFile(objectPath string, file io.Reader, fileSize int64) error {
	// Construire l'URL pour téléverser l'objet
	url := fmt.Sprintf("%s/%s", s.APIURL, objectPath)

	// Créer une requête PUT pour téléverser le fichier
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return fmt.Errorf("échec de la création de la requête : %v", err)
	}

	// Ajouter les en-têtes requis
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Amz-Decoded-Content-Length", fmt.Sprintf("%d", fileSize))

	// Envoyer la requête
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("échec de l'upload : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload échoué, statut : %s", resp.Status)
	}

	return nil
}

