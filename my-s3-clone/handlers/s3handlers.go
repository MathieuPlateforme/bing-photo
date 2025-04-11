package handlers

import (
    "io"
    "my-s3-clone/storage"
    "my-s3-clone/dto"
    "net/http"
    "github.com/gorilla/mux"
    "log"
    "time"
    "encoding/xml"
    "fmt"
    "os"
    "strconv"
    "errors"
)

// List all buckets
func HandleListBuckets(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        log.Println("Calling storage.ListBuckets to get the list of buckets.")
        buckets := s.ListBuckets()

        if len(buckets) == 0 {
            log.Println("No buckets found in storage.")
        } else {
            log.Printf("Found %d buckets.", len(buckets))
        }

        var bucketList []dto.Bucket
        for _, bucketName := range buckets {
            log.Printf("Adding bucket: %s", bucketName)
            bucketList = append(bucketList, dto.Bucket{
                Name:         bucketName,
                CreationDate: time.Now(),
            })
        }

        response := dto.ListAllMyBucketsResult{
            Buckets: bucketList,
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)

        log.Println("Encoding response as XML and sending it.")
        if err := xml.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Erreur lors de l'encodage des buckets en XML", http.StatusInternalServerError)
            log.Printf("Erreur lors de l'encodage des buckets: %v", err)
        }
    }
}

// Create a bucket
func HandleCreateBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        if r.Method != "PUT" {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
            return
        }

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        // Vérification si le bucket existe déjà
        exists, err := s.CheckBucketExists(bucketName) 
        if err != nil {
            http.Error(w, "Erreur lors de la vérification du bucket", http.StatusInternalServerError)
            return
        }

        if exists {
            http.Error(w, fmt.Sprintf("Bucket '%s' already exists", bucketName), http.StatusConflict)
            return
        }

        // Création du bucket si il n'existe pas
        err = s.CreateBucket(bucketName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Réponse pour indiquer que le bucket a été créé avec succès
        bucketResponse := dto.ListAllMyBucketsResult{
            Buckets: []dto.Bucket{
                {
                    Name:         bucketName,
                },
            },
        }

        w.Header().Set("Content-Type", "application/xml")
        w.Header().Set("Location", r.URL.String())
        w.WriteHeader(http.StatusOK)
        if err := xml.NewEncoder(w).Encode(bucketResponse); err != nil {
            http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
        }
    }
}

// Get bucket info or location
func HandleGetBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        log.Printf("Requête GET pour le bucket: %s", bucketName)

        // Gérer le paramètre de localisation
        locationParam := r.URL.Query().Get("location")
        log.Printf("Location Param: %s", locationParam)

        if locationParam != "" {
            log.Printf("Demande de localisation pour le bucket: %s", bucketName)
            w.Header().Set("Content-Type", "application/xml")
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`<LocationConstraint>us-east-1</LocationConstraint>`))
            return
        }

        // Vérifier si le bucket existe
        exists, err := s.CheckBucketExists(bucketName)
        if err != nil {
            log.Printf("Erreur lors de la vérification du bucket: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        if !exists {
            log.Printf("Bucket non trouvé: %s", bucketName)
            http.Error(w, fmt.Sprintf("Bucket '%s' not found", bucketName), http.StatusNotFound)
            return
        }

        // Si le bucket existe
        log.Printf("Bucket %s existe et est accessible", bucketName)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Bucket '%s' exists and is accessible.", bucketName)))
    }
}

// Add an object
func HandleAddObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        if bucketName == "" || objectName == "" {
            http.Error(w, "Bucket name and object name are required", http.StatusBadRequest)
            log.Printf("Bucket name or object name missing: bucketName=%s, objectName=%s", bucketName, objectName)
            return
        }

        log.Printf("Uploading object: %s to bucket: %s", objectName, bucketName)

        // Get the total content length from the X-Amz-Decoded-Content-Length header
        contentLength := r.Header.Get("X-Amz-Decoded-Content-Length")
        if contentLength == "" {
            log.Printf("Missing X-Amz-Decoded-Content-Length header")
            http.Error(w, "Missing X-Amz-Decoded-Content-Length header", http.StatusBadRequest)
            return
        }

        log.Printf("Total upload size: %s bytes", contentLength)

        // Process the uploaded object
        err := s.AddObject(bucketName, objectName, r.Body, r.Header.Get("X-Amz-Content-Sha256"))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            log.Printf("Error uploading object: %v", err)
            return
        }

        // Generate an ETag for the object
        eTag := `"1b2cf535f27731c974343645a3985328"`

        // Set the appropriate headers
        w.Header().Set("ETag", eTag)
        w.Header().Set("x-amz-id-2", "LriYPLdmOdAiIfgSm/F1YsViT1LW94/xUQxMsF7xiEb1a0wiIOIxl+zbwZ163pt7")
        w.Header().Set("x-amz-request-id", "0A49CE4060975EAC")
        w.Header().Set("Date", time.Now().Format(http.TimeFormat))

        // Send the response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte{}) // Empty body, as per S3 standard response

        // Log response status and body
        log.Printf("Response status: %d", http.StatusOK)
        log.Printf("Successfully uploaded object: %s in bucket: %s", objectName, bucketName)
    }
}

// Check if an object exists
func HandleCheckObjectExist(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        if bucketName == "" || objectName == "" {
            http.Error(w, "Bucket name and object name are required", http.StatusBadRequest)
            return
        }

        exists, lastModified, size, err := s.CheckObjectExist(bucketName, objectName)
        if err != nil || !exists {
            if !exists {
                http.Error(w, "Object not found", http.StatusNotFound)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
        w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
        w.WriteHeader(http.StatusOK)
    }
}

// Download an object
func HandleDownloadObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        objectName := vars["objectName"]

        // Récupérer les données du fichier et ses métadonnées
        data, fileInfo, err := s.GetObject(bucketName, objectName)
        if err != nil {
            if os.IsNotExist(err) {
                http.Error(w, "File not found", http.StatusNotFound)
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Envoyer les métadonnées dans les en-têtes HTTP
        w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", objectName))
        w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
        w.Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))

        // Envoyer le contenu du fichier
        w.Header().Set("Content-Type", "application/octet-stream")
        w.WriteHeader(http.StatusOK)

        if _, err := w.Write(data); err != nil {
            http.Error(w, "Failed to write file content", http.StatusInternalServerError)
        }
    }
}

// List objects in a bucket
func HandleListObjects(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        queryParams := r.URL.Query()
        prefix := queryParams.Get("prefix")
        marker := queryParams.Get("marker")
        maxKeys := queryParams.Get("max-keys")

        if maxKeys == "" {
            maxKeys = "1000"
        }

        maxKeysInt, err := strconv.Atoi(maxKeys)
        if err != nil {
            http.Error(w, "Invalid max-keys value", http.StatusBadRequest)
            return
        }

        objects, err := s.ListObjects(bucketName, prefix, marker, maxKeysInt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        if err := xml.NewEncoder(w).Encode(objects); err != nil {
            http.Error(w, "Erreur lors de l'encodage XML", http.StatusInternalServerError)
        }
    }
}

// Delete a bucket
func HandleDeleteBucket(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
        
        if bucketName == "" {
            http.Error(w, "Bucket name is required", http.StatusBadRequest)
            return
        }

        // Tenter de supprimer le bucket
        err := s.DeleteBucket(bucketName)
        if err != nil {
            // Si l'erreur indique que le bucket n'existe pas, renvoyer un code 404
            if os.IsNotExist(err) {
                log.Printf("Bucket %s does not exist", bucketName)
                http.Error(w, fmt.Sprintf("Bucket %s does not exist", bucketName), http.StatusNotFound)
                return
            }
            // Pour toute autre erreur, renvoyer un code 500
            log.Printf("Error deleting bucket %s: %v", bucketName, err)
            http.Error(w, "Failed to delete bucket", http.StatusInternalServerError)
            return
        }

        // Répondre avec succès si le bucket est supprimé
        log.Printf("Bucket %s deleted successfully", bucketName)
        w.WriteHeader(http.StatusNoContent)
    }
}


// Batch delete objects
func HandleDeleteObject(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }
        log.Printf("Received POST ?delete request for batch deletion: %s %s", r.Method, r.URL.Path)

        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading request body", http.StatusInternalServerError)
            log.Printf("Error reading request body: %v", err)
            return
        }
        log.Printf("Request body: %s", string(body))

        var deleteReq dto.DeleteObjectRequest
        err = xml.Unmarshal(body, &deleteReq)
        if err != nil {
            http.Error(w, "Error parsing XML", http.StatusBadRequest)
            log.Printf("Error parsing XML: %v", err)
            return
        }

        var deletedObjects []dto.Deleted
        for _, objectToDelete := range deleteReq.Objects {
            log.Printf("Attempting to delete object: %s", objectToDelete.Key)
            err := s.DeleteObject(bucketName, objectToDelete.Key)
            if err != nil {
                if errors.Is(err, os.ErrNotExist) { // Vérifie si l'erreur correspond à l'objet non trouvé
                    http.Error(w, "Object not found", http.StatusNotFound)
                    log.Printf("Object not found: %s", objectToDelete.Key)
                    continue 
                }
                http.Error(w, "Error deleting object", http.StatusInternalServerError)
                log.Printf("Error deleting object %s: %v", objectToDelete.Key, err)
                return
            }
            log.Printf("Successfully deleted object: %s", objectToDelete.Key)

            deletedObjects = append(deletedObjects, dto.Deleted{Key: objectToDelete.Key})
        }

        deleteResult := dto.DeleteResult{
            DeletedResult: deletedObjects,
        }

        response, err := xml.Marshal(deleteResult)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            log.Printf("Error generating XML response: %v", err)
            return
        }

        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)

                // Log response status and body
        log.Printf("Response status: %d", http.StatusOK)
        log.Printf("Response body: %s", string(response))
    }
}

func HandleBucketLocation(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]
            // La chaîne représentant la date
    dateString := "2024-09-16T10:12:24Z"

    // Convertir la chaîne en time.Time
    creationDate, err := time.Parse(time.RFC3339, dateString)
    if err != nil {
        fmt.Printf("Error parsing date: %v\n", err)
        return
    }
        bucket := dto.Bucket{
            Name:         bucketName,
            CreationDate: creationDate,
            LocationConstraint : "us-east-1",
        }

        response, err := xml.Marshal(bucket)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            log.Printf("Error generating XML response: %v", err)
            return
        }

        
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
            
    }
}

func HandleBucketLockConfig(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        dateString := "2024-09-16T10:12:24Z"

        // Convertir la chaîne en time.Time
        creationDate, err := time.Parse(time.RFC3339, dateString)
        if err != nil {
            fmt.Printf("Error parsing date: %v\n", err)
            return
        }

        bucket := dto.Bucket{
            Name:         bucketName,
            CreationDate: creationDate,
            ObjectLockConfig : "true",
        }
    

        response, err := xml.Marshal(bucket)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            log.Printf("Error generating XML response: %v", err)
            return
        }

        
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
            
            
    }
}


func HandleBucketDelimiter(s storage.Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        bucketName := vars["bucketName"]

        dateString := "2024-09-16T10:12:24Z"

        // Convertir la chaîne en time.Time
        creationDate, err := time.Parse(time.RFC3339, dateString)
        if err != nil {
            fmt.Printf("Error parsing date: %v\n", err)
            return
        }

        bucket := dto.Bucket{
            Name:         bucketName,
            CreationDate: creationDate,
            ObjectDelimiter : "true",
        }
    

        response, err := xml.Marshal(bucket)
        if err != nil {
            http.Error(w, "Error generating XML response", http.StatusInternalServerError)
            log.Printf("Error generating XML response: %v", err)
            return
        }

        
        w.Header().Set("Content-Type", "application/xml")
        w.WriteHeader(http.StatusOK)
        w.Write(response)
            
            
    }
}

type MoveObjectRequest struct {
	XMLName   xml.Name        `xml:"Move"`
	Objects   []ObjectToMove  `xml:"Object"`
	TargetBucket string       `xml:"TargetBucket"`
}

// ObjectToMove représente un objet à déplacer
type ObjectToMove struct {
	Key string `xml:"Key"`
}

// HandleMoveObject gère les requêtes pour déplacer des objets entre des buckets
func HandleMoveObject(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		log.Printf("Received POST ?move request for moving objects: %s %s", r.Method, r.URL.Path)

		vars := mux.Vars(r)
		sourceBucket := vars["bucketName"]

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			log.Printf("Error reading request body: %v", err)
			return
		}
		log.Printf("Request body: %s", string(body))

		var moveReq MoveObjectRequest
		err = xml.Unmarshal(body, &moveReq)
		if err != nil {
			http.Error(w, "Error parsing XML", http.StatusBadRequest)
			log.Printf("Error parsing XML: %v", err)
			return
		}

		if moveReq.TargetBucket == "" {
			http.Error(w, "TargetBucket is missing in the request", http.StatusBadRequest)
			log.Println("Error: TargetBucket is missing")
			return
		}

		var movedObjects []dto.Deleted
		for _, objectToMove := range moveReq.Objects {
			log.Printf("Attempting to move object: %s", objectToMove.Key)

			// Copier l'objet
			err := s.CopyObject(sourceBucket, objectToMove.Key, moveReq.TargetBucket, objectToMove.Key)
			if err != nil {
				http.Error(w, "Error moving object", http.StatusInternalServerError)
				log.Printf("Error moving object %s: %v", objectToMove.Key, err)
				return
			}

			// Supprimer l'objet source
			err = s.DeleteObject(sourceBucket, objectToMove.Key)
			if err != nil {
				http.Error(w, "Error deleting source object after move", http.StatusInternalServerError)
				log.Printf("Error deleting object %s: %v", objectToMove.Key, err)
				return
			}

			log.Printf("Successfully moved object: %s", objectToMove.Key)
			movedObjects = append(movedObjects, dto.Deleted{Key: objectToMove.Key})
		}

		moveResult := dto.DeleteResult{
			DeletedResult: movedObjects,
		}

		response, err := xml.Marshal(moveResult)
		if err != nil {
			http.Error(w, "Error generating XML response", http.StatusInternalServerError)
			log.Printf("Error generating XML response: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

		// Log response status and body
		log.Printf("Response status: %d", http.StatusOK)
		log.Printf("Response body: %s", string(response))
	}
}