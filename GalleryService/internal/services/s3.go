package services

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	Name               string    `xml:"Name"`
	CreationDate       time.Time `xml:"CreationDate"`
	LocationConstraint string    `xml:"LocationConstraint,omitempty"`
	ObjectLockConfig   string    `xml:"ObjectLockConfiguration,omitempty"`
	ObjectDelimiter    string    `xml:"ObjectDelimiter,omitempty"`
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

func (s *S3Service) MoveObject(sourceBucket, sourceKey, targetBucket string) error {
	// Construire l'URL pour déplacer l'objet
	url := fmt.Sprintf("%s/%s/?move", s.APIURL, sourceBucket)

	// Construire le corps de la requête en XML
	payload := fmt.Sprintf(`
	<Move>
		<Object>
			<Key>%s</Key>
		</Object>
		<TargetBucket>%s</TargetBucket>
	</Move>

	`, sourceKey, targetBucket)

	// Créer la requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return fmt.Errorf("échec de la création de la requête de déplacement : %v", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	// Envoyer la requête
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("échec de la requête de déplacement : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de réponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("échec de la requête de déplacement, statut : %s", resp.Status)
	}

	// Déplacement réussi
	return nil
}

func (s *S3Service) DownloadFile(path string, w io.Writer) error {
	// Construire l'URL pour télécharger le fichier
	url := fmt.Sprintf("%s/%s", s.APIURL, path)

	// Envoyer une requête HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("échec de la requête de téléchargement : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de réponse
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("échec du téléchargement, statut : %s", resp.Status)
	}

	// Copier le contenu de la réponse dans le writer fourni
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("échec de l'écriture des données téléchargées : %v", err)
	}

	log.Printf("Fichier téléchargé avec succès depuis %s", url)
	return nil
}

func (s *S3Service) DeleteObject(bucketName, objectName string) error {
	// Journaliser la tentative de suppression
	log.Printf("Tentative de suppression : bucket=%s, key=%s", bucketName, objectName)

	// Construire l'URL pour supprimer l'objet
	url := fmt.Sprintf("%s/%s/?delete=", s.APIURL, bucketName)
	log.Printf("URL construite pour la suppression : %s", url)

	// Construire le corps de la requête en XML
	payload := fmt.Sprintf(`
    <Delete>
        <Object>
            <Key>%s</Key>
        </Object>
    </Delete>
    `, objectName)
	log.Printf("Corps XML généré : %s", payload)

	// Créer une requête HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return fmt.Errorf("échec de la création de la requête de suppression : %v", err)
	}
	req.Header.Set("Content-Type", "application/xml")

	// Journaliser les en-têtes
	log.Printf("Headers envoyés : %v", req.Header)

	// Envoyer la requête HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("échec de la requête de suppression : %v", err)
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse pour le log
	responseBody, _ := io.ReadAll(resp.Body)
	log.Printf("Code de réponse HTTP : %d, Corps : %s", resp.StatusCode, string(responseBody))

	// Vérifier le code de réponse
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("échec de la suppression de l'objet, statut : %s", resp.Status)
	}

	// Journaliser la réussite
	log.Printf("Objet supprimé avec succès : %s/%s", bucketName, objectName)
	return nil
}

func (s *S3Service) GetFilesInAlbum(bucketName string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/", s.APIURL, bucketName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("échec de la récupération des fichiers pour l'album %s : %v", bucketName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("échec de la requête, statut HTTP : %d", resp.StatusCode)
	}

	var listResponse struct {
		Objects []struct {
			Key string `xml:"Key"`
		} `xml:"Contents"`
	}

	err = xml.NewDecoder(resp.Body).Decode(&listResponse)
	if err != nil {
		return nil, fmt.Errorf("échec du décodage XML : %v", err)
	}

	fileNames := []string{}
	for _, obj := range listResponse.Objects {
		fileNames = append(fileNames, obj.Key)
	}

	log.Printf("Fichiers récupérés pour l'album %s : %v", bucketName, fileNames)
	return fileNames, nil
}

func (s *S3Service) DownloadTempFile(bucketName, objectName string) (string, error) {
	localPath := fmt.Sprintf("/tmp/%s", objectName)

	file, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("échec de la création du fichier temporaire: %v", err)
	}
	defer file.Close()

	filePath := fmt.Sprintf("%s/%s", bucketName, objectName)
	err = s.DownloadFile(filePath, file)
	if err != nil {
		return "", fmt.Errorf("échec du téléchargement du fichier: %v", err)
	}

	log.Printf("Fichier téléchargé temporairement : %s", localPath)
	return localPath, nil
}
