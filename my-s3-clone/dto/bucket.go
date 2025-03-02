package dto

import (
    "encoding/xml"
    "time"
)

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
