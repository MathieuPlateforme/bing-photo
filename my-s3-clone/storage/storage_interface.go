package storage

import (
	"io"
	"time"
	"my-s3-clone/dto"

)

// Storage interface définissant les méthodes de gestion des objets et des buckets
type Storage interface {
    AddObject(bucketName, objectName string, data io.Reader, contentSha256 string) error
    DeleteObject(bucketName, objectName string) error
    DeleteBucket(bucketName string) error
    GetObject(bucketName, objectName string) ([]byte, dto.FileInfo, error)
    CheckObjectExist(bucketName, objectName string) (bool, time.Time, int64, error)
    CheckBucketExists(bucketName string) (bool, error)
    ListBuckets() []string
    ListObjects(bucketName, prefix, marker string, maxKeys int) (dto.ListObjectsResponse, error)
    CreateBucket(bucketName string) error
}


