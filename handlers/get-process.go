package handlers

import (
	"encoding/json"
	"goxl/database"
	"goxl/util"
	"net/http"
)

type Body struct {
	ProcessId string `json:"process_id"`
}

func GetProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.RespondWithError(w, http.StatusBadRequest, "invalid method")
		return
	}
	// Decode the JSON body
	var body Body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	// Check if ProcessId is empty
	if body.ProcessId == "" {
		util.RespondWithError(w, http.StatusBadRequest, "process_id is required")
		return
	}

	var db database.Database

	connErr := db.Connect()
	if connErr != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Something went giga wrong")
	}

	prc, err := db.GetProcess(body.ProcessId)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Process does not exist, or something else went wrong")

	}

	defer db.Disconnect()

	// Example of a successful response
	util.RespondWithJson(w, http.StatusOK, prc)
}
