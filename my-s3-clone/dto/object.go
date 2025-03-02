package dto

import (
    "encoding/xml"
    "time"
)

type ListObjectsResponse struct {
    XMLName     xml.Name `xml:"ListBucketResult"`
    Xmlns       string   `xml:"xmlns,attr"`
    Name        string   `xml:"Name"`
    Prefix      string   `xml:"Prefix"`
    Marker      string   `xml:"Marker"`
    MaxKeys     int      `xml:"MaxKeys"`
    IsTruncated bool     `xml:"IsTruncated"`
    Contents    []Object `xml:"Contents"`
}

type Object struct {
    Key          string    `xml:"Key"`
    LastModified time.Time `xml:"LastModified"`
    Size         int       `xml:"Size"`
}
