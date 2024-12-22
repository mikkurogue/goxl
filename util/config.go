package util

import (
	"encoding/json"
	"os"

	"github.com/fatih/color"
)

type Config struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	ColumnName   string `json:"column_name"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	Type         string `json:"type"`
}

const ConfigName = "goxl.config.json"

func GenerateDefaultConfig() {
	config := Config{
		Name: "MyProduct",
		Columns: []Column{
			{
				ColumnName:   "id",
				IsPrimaryKey: true,
				Type:         "integer",
			},
			{
				ColumnName:   "name",
				IsPrimaryKey: false,
				Type:         "text",
			},
			{
				ColumnName:   "vehicle",
				IsPrimaryKey: false,
				Type:         "text",
			},
		},
	}

	jsonData, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		color.Red("Could not marshal config: ", err)
		os.Exit(0)
	}

	fileName := ConfigName
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		color.Red("Could not write to file: ", err)
		os.Exit(0)
	}

	color.Green("Created initial config file.")
}

func HasConfig() bool {
	_, err := os.ReadFile(ConfigName)
	return !os.IsNotExist(err)
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigName)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
