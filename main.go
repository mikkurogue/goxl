package main

import (
	"fmt"
	"goxl/bootstrap"
	"goxl/database"
	"goxl/middleware"
	"goxl/util"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {
	if !util.HasConfig() {
		util.GenerateDefaultConfig()
		color.HiGreen("Default config generated. Please adjust the config to the requirements and restart the api.")
		os.Exit(0)
	}

	config, err := util.ReadConfig()
	if err != nil {
		color.HiRed("Could not read or find a valid configuration")
		panic(0)
	}

	var db database.Database
	db.Connect()
	defer db.Disconnect()

	db.CreateTable(config)

	mux := bootstrap.Bootstrap()

	// wrap mux with cors middleware
	handler := middleware.CORS(mux)

	fmt.Println("Server running on http://localhost:9292")

	error := http.ListenAndServe(":9292", handler)
	if error != nil {
		fmt.Println("Error starting server", error)
	}
}
