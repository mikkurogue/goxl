package bootstrap

import (
	"goxl/handlers"
	"log"
	"net/http"
	"os"
)

func Bootstrap() *http.ServeMux {
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatal("Could not create uploads dir")
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload-file", handlers.UploadFile)
	mux.HandleFunc("/api/get-process", handlers.GetProcessDetails)
	mux.HandleFunc("/api/download", handlers.DownloadFileFromProcess)

	return mux
}

func RemoveDbStore() {
	os.Remove("./store.db")
}
