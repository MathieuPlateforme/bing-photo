package dto

import (
	"os"
	"time"
)
// FileInfo représente les métadonnées d'un fichier (objet)
type FileInfo interface {
    Name() string       // Nom de base du fichier
    Size() int64        // Taille logique du fichier en octets
    Mode() os.FileMode  // Informations sur le mode de fichier
    ModTime() time.Time // Heure de dernière modification
    IsDir() bool        // Indique si c'est un répertoire
    Sys() interface{}   // Données spécifiques au système sous-jacent
}

// FileInfoWrapper encapsule un os.FileInfo pour implémenter l'interface FileInfo
type FileInfoWrapper struct {
	FileInfo os.FileInfo
}

// Implémentation des méthodes de l'interface FileInfo

func (fi *FileInfoWrapper) Name() string {
	return fi.FileInfo.Name()
}

func (fi *FileInfoWrapper) Size() int64 {
	return fi.FileInfo.Size()
}

func (fi *FileInfoWrapper) Mode() os.FileMode {
	return fi.FileInfo.Mode()
}

func (fi *FileInfoWrapper) ModTime() time.Time {
	return fi.FileInfo.ModTime()
}

func (fi *FileInfoWrapper) IsDir() bool {
	return fi.FileInfo.IsDir()
}

func (fi *FileInfoWrapper) Sys() interface{} {
	return fi.FileInfo.Sys()
}
