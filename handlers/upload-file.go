package handlers

import (
	"fmt"
	"goxl/database"
	"goxl/util"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// UploadFile handles file uploads via multipart/form-data
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Parse the multipart form with a max memory limit (e.g., 10 MB)
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		fmt.Println("Error parsing form:", err)
		util.RespondWithError(w, http.StatusBadRequest, "Error parsing form: "+err.Error())
		return
	}

	name := r.FormValue("name")
	if name == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing 'name' field")
		return
	}

	// Retrieve the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving file:", err)
		util.RespondWithError(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// fmt.Println("Uploaded File Name:", header.Filename)
	// fmt.Println("Uploaded File Size:", header.Size)
	// fmt.Println("Uploaded File Header:", header.Header)

	var db database.Database
	config, _ := util.ReadConfig()

	connErr := db.Connect()
	if connErr != nil {
		color.Red(connErr.Error())
	}

	newUuid, err := uuid.NewUUID()
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Something went wrong on the server, contact support")
		return
	}

	// Create the subfolder path using the UUID
	subfolderPath := filepath.Join("uploads", newUuid.String())

	// Ensure the folder exists
	err = os.MkdirAll(subfolderPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating subfolder:", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Error creating subfolder: "+err.Error())
		return
	}

	// Define the destination file path inside the subfolder
	dstPath := filepath.Join(subfolderPath, header.Filename) // You can change "uploaded_file" to the desired file name
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println("Error creating file on server:", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Error creating file on server: "+err.Error())
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error saving file:", err)
		util.RespondWithError(w, http.StatusInternalServerError, "Error saving file: "+err.Error())
		return
	}

	var cols []string
	for _, v := range config.Columns {
		cols = append(cols, v.ColumnName)
	}

	var prc = database.ProcessRow{
		Id:       newUuid.String(),
		FileName: header.Filename,
		FileSize: header.Size,
	}

	insErr := db.AddUpload(prc)
	if insErr != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Something went wrong with creating upload instance...")
	}

	db.Disconnect()

	response := map[string]string{
		"message": "File uploaded successfully",
		"name":    name,
		"file":    header.Filename,
	}
	util.RespondWithJson(w, http.StatusOK, response)
}
