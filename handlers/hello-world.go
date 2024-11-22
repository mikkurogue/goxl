package handlers

import (
	"goxl/util"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJson(w, http.StatusOK, map[string]string{"hello": "world"})
}
