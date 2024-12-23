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
)

// UploadFile handles file uploads via multipart/form-data
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Debug: Print the request method and headers
	fmt.Println("Method:", r.Method)
	fmt.Println("Headers:", r.Header)

	// Ensure the method is POST
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

	// Debug: Print parsed form values
	fmt.Println("Form Values:", r.Form)
	fmt.Println("Form Files:", r.MultipartForm.File)

	// Retrieve the 'name' field
	name := r.FormValue("name")
	if name == "" {
		util.RespondWithError(w, http.StatusBadRequest, "Missing 'name' field")
		return
	}

	// Retrieve the 'extension' field (optional)
	extension := r.FormValue("extension")

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

	dstPath := filepath.Join("uploads")
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

	insErr := db.InsertRow(
		cols,
		[]string{
			"0",
			"Mustang",
			"Ford Mustang GT",
		},
	)
	if insErr != nil {
		color.Red(insErr.Error())
	}

	db.Disconnect()

	response := map[string]string{
		"message":   "File uploaded successfully",
		"name":      name,
		"file":      header.Filename,
		"extension": extension,
	}
	util.RespondWithJson(w, http.StatusOK, response)
}
