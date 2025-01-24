package storage

import (
    "strings"
    "os"
    "path/filepath"
    "log"
    "fmt"
    "io"
    "bufio"  
    "strconv"
    "time"
    "my-s3-clone/dto"
)

// FileStorage implémente l'interface Storage avec un stockage basé sur le système de fichiers
type FileStorage struct{}

const storageRoot = "/mydata/data"

func ProcessChunkedStream(reader io.Reader, writer io.Writer) error {
    bufReader := bufio.NewReader(reader)
    log.Println("Started processing chunked stream")

    totalBytesProcessed := int64(0)

    for {
        // Read chunk size line
        line, err := bufReader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading chunk size: %v", err)
            return fmt.Errorf("error reading chunk size: %v", err)
        }
        log.Printf("Received chunk size line: %s", line)

        // Parse chunk size (ignore chunk extensions)
        line = strings.TrimSpace(line)
        parts := strings.SplitN(line, ";", 2)
        chunkSizeHex := parts[0]

        chunkSize, err := strconv.ParseInt(chunkSizeHex, 16, 64)
        if err != nil {
            log.Printf("Error parsing chunk size: %v", err)
            return fmt.Errorf("error parsing chunk size: %v", err)
        }

        log.Printf("Parsed chunk size: %d", chunkSize)

        // End of stream (zero-size chunk)
        if chunkSize == 0 {
            log.Println("Received final chunk (size 0), finishing")
            break
        }

        // Copy chunk data to writer
        if _, err := io.CopyN(writer, bufReader, chunkSize); err != nil {
            log.Printf("Error reading chunk data: %v", err)
            return fmt.Errorf("error reading chunk data: %v", err)
        }

        totalBytesProcessed += chunkSize
        log.Printf("Successfully processed a chunk, total bytes processed: %d", totalBytesProcessed)

        // Discard the CRLF after the chunk
        if _, err := bufReader.Discard(2); err != nil {
            log.Printf("Error discarding CRLF: %v", err)
            return fmt.Errorf("error discarding CRLF: %v", err)
        }

        // Log chunk signature 
        if len(parts) > 1 {
            chunkSignature := parts[1]
            log.Printf("Chunk signature: %s", chunkSignature)
        }
    }

    log.Printf("Completed processing chunked stream, total bytes processed: %d", totalBytesProcessed)
    return nil
}


// Ajout d'un objet dans un bucket
func (fs *FileStorage) AddObject(bucketName, objectName string, data io.Reader, contentSha256 string) error {
    log.Printf("Starting object upload: %s in bucket: %s", objectName, bucketName)

    objectPath, err := getUniqueObjectPath(bucketName, objectName)
    if err != nil {
        log.Printf("Failed to create object path for %s in bucket %s: %v", objectName, bucketName, err)
        return fmt.Errorf("Failed to create object path: %v", err)
    }

    log.Printf("Object path created: %s", objectPath)

    file, err := os.Create(objectPath)
    if err != nil {
        log.Printf("Failed to create file: %s, error: %v", objectPath, err)
        return fmt.Errorf("Failed to create file: %v", err)
    }
    defer file.Close()

    log.Printf("Writing data to object: %s", objectPath)

    if err := writeObjectToFile(data, file, contentSha256); err != nil {
        log.Printf("Error writing object to file: %v", err)
        return err
    }

    log.Printf("Successfully uploaded file: %s", objectPath)
    return nil
}

// Fonction pour obtenir un chemin unique si l'objet existe déjà
func getUniqueObjectPath(bucketName, objectName string) (string, error) {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)
    if _, err := os.Stat(objectPath); os.IsNotExist(err) {
        return objectPath, nil
    }

    log.Printf("Object already exists, generating new name for: %s", objectPath)

    objectNameWithoutExt := strings.TrimSuffix(objectName, filepath.Ext(objectName))
    extension := filepath.Ext(objectName)
    newObjectName := objectNameWithoutExt
    suffix := 1

    for {
        newObjectName = fmt.Sprintf("%s-%d%s", objectNameWithoutExt, suffix, extension)
        newObjectPath := filepath.Join(storageRoot, bucketName, newObjectName)
        if _, err := os.Stat(newObjectPath); os.IsNotExist(err) {
            return newObjectPath, nil
        }
        suffix++
    }
}

// Fonction qui gère l'écriture du flux dans le fichier
func writeObjectToFile(data io.Reader, file *os.File, contentSha256 string) error {
    if contentSha256 == "STREAMING-AWS4-HMAC-SHA256-PAYLOAD" {
        log.Println("Processing as chunked stream")
        if err := ProcessChunkedStream(data, file); err != nil {
            log.Printf("Failed to write chunked data: %v", err)
            return fmt.Errorf("Failed to write chunked data: %v", err)
        }
    } else {
        log.Println("Processing as regular stream")
        if _, err := io.Copy(file, data); err != nil {
            log.Printf("Failed to write data: %v", err)
            return fmt.Errorf("Failed to write data: %v", err)
        }
    }
    return nil
}

// Lister les objets dans un bucket
func (fs *FileStorage) ListObjects(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error) {
    bucketPath := filepath.Join(storageRoot, bucketName)

    objects, err := filepath.Glob(filepath.Join(bucketPath, prefix+"*"))
    if err != nil {
        return dto.ListObjectsResponse{}, fmt.Errorf("error while listing objects: %v", err)
    }

    response := dto.ListObjectsResponse{
        Xmlns:       "http://s3.amazonaws.com/doc/2006-03-01/",
        Name:        bucketName,
        Prefix:      prefix,
        Marker:      marker,
        MaxKeys:     maxKeys,
        IsTruncated: false,
        Contents:    make([]dto.Object, 0),
    }

    for i, object := range objects {
        if i >= maxKeys {
            response.IsTruncated = true
            break
        }

        fileInfo, err := os.Stat(object)
        if err != nil {
            return dto.ListObjectsResponse{}, fmt.Errorf("error retrieving file info: %v", err)
        }

        response.Contents = append(response.Contents, dto.Object{
            Key:          filepath.Base(object),
            LastModified: fileInfo.ModTime(),
            Size:         int(fileInfo.Size()),
        })
    }

    return response, nil
}

// Lister les buckets
func (fs *FileStorage) ListBuckets() []string {
    var buckets []string

    // Ajout de log pour vérifier si le répertoire existe
    log.Printf("Vérification de l'existence du répertoire de stockage des buckets : %s", storageRoot)

    files, err := os.ReadDir(storageRoot)
    if err != nil {
        log.Printf("Erreur lors de la lecture du répertoire %s : %v", storageRoot, err)
        return buckets
    }

    // Ajout de log pour voir combien de fichiers/répertoires sont trouvés
    log.Printf("Nombre d'éléments trouvés dans le répertoire %s : %d", storageRoot, len(files))

    // Parcourir chaque élément trouvé
    for _, file := range files {
        if file.IsDir() {
            // Ajout de log pour chaque répertoire trouvé
            log.Printf("Bucket trouvé : %s", file.Name())
            buckets = append(buckets, file.Name())
        } else {
            // Si un fichier est trouvé au lieu d'un répertoire, le signaler
            log.Printf("Fichier non attendu trouvé dans le répertoire des buckets : %s", file.Name())
        }
    }

    // Ajout de log pour indiquer la fin de la fonction et le nombre total de buckets trouvés
    log.Printf("Nombre total de buckets trouvés : %d", len(buckets))

    return buckets
}

// Créer un bucket
func (fs *FileStorage) CreateBucket(bucketName string) error {
    bucketPath := filepath.Join(storageRoot, bucketName)
    if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
        return err
    }
    return nil
}

// Récupération d'un objet dans un bucket
func (fs *FileStorage) GetObject(bucketName, objectName string) ([]byte, dto.FileInfo, error) {
	objectPath := filepath.Join(storageRoot, bucketName, objectName)
	log.Printf("Tentative de récupération de l'objet : %s", objectPath)

	// Lire le fichier
	data, err := os.ReadFile(objectPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture de l'objet: %v", err)
		return nil, nil, err
	}

	// Récupérer les métadonnées du fichier
	FileInfo, err := os.Stat(objectPath)
	if err != nil {
		log.Printf("Erreur lors de la récupération des métadonnées du fichier: %v", err)
		return nil, nil, err
	}

	// Retourner le contenu du fichier et les métadonnées encapsulées dans fileInfoWrapper
	return data, &dto.FileInfoWrapper{FileInfo: FileInfo}, nil
}

// Vérification de l'existence d'un objet dans un bucket
func (fs *FileStorage) CheckObjectExist(bucketName, objectName string) (bool, time.Time, int64, error) {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)

    fileInfo, err := os.Stat(objectPath)
    if os.IsNotExist(err) {
        return false, time.Time{}, 0, nil
    } else if err != nil {
        log.Printf("Error checking object: %v", err)
        return false, time.Time{}, 0, fmt.Errorf("error checking object existence: %v", err)
    }

    return true, fileInfo.ModTime(), fileInfo.Size(), nil
}

// Vérification de l'existence d'un bucket
func (fs *FileStorage) CheckBucketExists(bucketName string) (bool, error) {
    bucketPath := filepath.Join(storageRoot, bucketName)
    if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
        return false, nil
    } else if err != nil {
        return false, err
    }
    return true, nil
}

// Suppression d'un bucket
func (fs *FileStorage) DeleteBucket(bucketName string) error {
    bucketPath := filepath.Join(storageRoot, bucketName)

    if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
        log.Printf("Bucket %s does not exist", bucketName)
        return err
    }

    err := os.RemoveAll(bucketPath)
    if err != nil {
        log.Printf("Failed to delete bucket %s: %v", bucketName, err)
        return err
    }

    log.Printf("Bucket %s successfully deleted", bucketName)
    return nil
}

// Suppression d'un objet dans un bucket
func (fs *FileStorage) DeleteObject(bucketName, objectName string) error {
    objectPath := filepath.Join(storageRoot, bucketName, objectName)

    if _, err := os.Stat(objectPath); os.IsNotExist(err) {
        log.Printf("Object %s does not exist in bucket %s", objectName, bucketName)
        return fmt.Errorf("object not found: %w", err) // Retourne une erreur "object not found" encapsulant l'erreur 404
    }

    err := os.Remove(objectPath)
    if err != nil {
        log.Printf("Failed to delete object %s in bucket %s: %v", objectName, bucketName, err)
        return err
    }

    log.Printf("Object %s in bucket %s successfully deleted", objectName, bucketName)
    return nil
}

func (fs *FileStorage) CopyObject(sourceBucket, sourceKey, targetBucket, targetKey string) error {
	sourcePath := filepath.Join(storageRoot, sourceBucket, sourceKey)
	targetPath := filepath.Join(storageRoot, targetBucket, targetKey)

	// Assurez-vous que le répertoire cible existe
	if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return fmt.Errorf("impossible de créer le répertoire cible : %v", err)
	}

	// Copier le fichier
	input, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("impossible d'ouvrir le fichier source : %v", err)
	}
	defer input.Close()

	output, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier cible : %v", err)
	}
	defer output.Close()

	if _, err := io.Copy(output, input); err != nil {
		return fmt.Errorf("erreur lors de la copie : %v", err)
	}

	return nil
}
