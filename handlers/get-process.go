package handlers

import (
	"archive/zip"
	"encoding/json"
	"goxl/database"
	"goxl/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Body struct {
	ProcessId string `json:"process_id"`
}

func GetProcessDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.RespondWithError(w, http.StatusBadRequest, "invalid method")
		return
	}
	// Decode the JSON body
	var body Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid JSON body",
		)
		return
	}
	defer r.Body.Close()

	// Check if ProcessId is empty
	if body.ProcessId == "" {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"process_id is required",
		)
		return
	}

	var db database.Database

	connErr := db.Connect()
	if connErr != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Something went giga wrong",
		)
	}

	prc, err := db.GetProcess(body.ProcessId)
	if err != nil {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"Process does not exist, or something else went wrong",
		)

	}

	defer db.Disconnect()

	// Example of a successful response
	util.RespondWithJson(w, http.StatusOK, prc)
}

func DownloadFileFromProcess(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		util.RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid method",
		)
		return
	}

	var body Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		return
	}
	defer r.Body.Close()

	processDir := filepath.Join("uploads", body.ProcessId)
	tempDir := filepath.Join("uploads", "temp")
	archiveFile := filepath.Join(tempDir, body.ProcessId+".zip")

	// Check if the process directory exists
	if _, err := os.Stat(processDir); os.IsNotExist(err) {
		util.RespondWithError(
			w,
			http.StatusNotFound,
			"folder for process does not exist",
		)
		return
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to create temp directory",
		)
		return
	}

	// Create a zip archive of the directory
	if err := createZipFromFolder(processDir, archiveFile); err != nil {
		util.RespondWithError(
			w,
			http.StatusInternalServerError,
			"failed to create archive: "+err.Error(),
		)
		return
	}
	defer os.Remove(archiveFile) // Clean up the temporary file after serving

	// Serve the file for download
	w.Header().Set("Content-Disposition", "attachment; filename="+body.ProcessId+".zip")
	w.Header().Set("Content-Type", "application/zip")
	http.ServeFile(w, r, archiveFile)
}

func createZipFromFolder(src, dst string) error {

	zipFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {

			return err
		}

		if info.IsDir() {

			return nil // skip directories
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {

			return err
		}

		file, err := os.Open(path)
		if err != nil {

			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)

		return err
	})

	return nil
}
