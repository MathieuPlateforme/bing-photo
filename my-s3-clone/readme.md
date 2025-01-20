# My S3 Clone

My S3 Clone est une API légère qui reproduit les fonctionnalités de base d'un service de stockage de type S3 en utilisant le système de fichier comme stockage. Elle permet de créer un alias de son server à l'aide du CLI minIO, créer des buckets, télécharger et récupérer des fichiers, lister les fichiers présents dans un bucket et supprimer des fichiers en suivant le protocole S3.

## Fonctionnalités

- **Créer un Bucket** : Crée un bucket de stockage dans MinIO.
- **Uploader un Objet** : Télécharge un objet dans un bucket.
- **Lister les Buckets** : Récupère la liste de tous les buckets.
- **Récupérer un Objet** : Récupère un objet spécifique depuis un bucket.
- **Supprimer un Objet** : Supprime un objet d'un bucket.
- **Supprimer un Bucket** : Supprime un bucket de MinIO.

## Prérequis

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Installation et Lancement

1. Clonez le dépôt :

    ```bash
    git clone https://github.com/votre-utilisateur/my-s3-clone.git
    cd my-s3-clone
    ```

2. Construisez et démarrez les conteneurs avec Docker Compose :

    ```bash
    docker-compose up --build
    ```

test


