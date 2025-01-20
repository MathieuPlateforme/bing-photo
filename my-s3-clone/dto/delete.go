package dto

import (
    "encoding/xml"
)
// type DeleteObjectRequest struct {
// 	Quiet  bool `xml:"Quiet"`
// 	Object struct {
// 		Key string `xml:"Key"`
// 	} `xml:"Object"`
// }

type DeleteResult struct {
	DeletedResult []Deleted `xml:"Deleted"`
}

type Deleted struct {
	Key string `xml:"Key"`
}

// DeleteObjectRequest représente la requête de suppression d'objets en batch
type DeleteObjectRequest struct {
    XMLName xml.Name         `xml:"Delete"`
    Objects []ObjectToDelete  `xml:"Object"`
}

// ObjectToDelete représente un objet à supprimer
type ObjectToDelete struct {
    Key string `xml:"Key"`
}
