package main

import (
    "log"
    "net/http"
    "os"
    "my-s3-clone/router"
)

func main() {
    if _, err := os.Stat("./buckets"); os.IsNotExist(err) {
        log.Printf("Le répertoire 'buckets' n'existe pas. Création...")
        if err := os.Mkdir("./buckets", os.ModePerm); err != nil {
            log.Fatalf("Erreur lors de la création du répertoire 'buckets': %v", err)
        }
    }

    r := router.SetupRouter() 
    log.Println("Serving on :9090")
    log.Fatal(http.ListenAndServe(":9090", r))
}
