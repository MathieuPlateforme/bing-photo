package dto

import (
    "encoding/xml"
)

type ErrorResponse struct {
    XMLName   xml.Name `xml:"Error"`
    Code      string   `xml:"Code"`
    Message   string   `xml:"Message"`
    BucketName string  `xml:"BucketName,omitempty"`
    RequestId string   `xml:"RequestId"`
    HostId    string   `xml:"HostId"`
}
