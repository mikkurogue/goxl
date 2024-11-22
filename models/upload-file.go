package models

import "mime/multipart"

type UploadFile struct {
	Name string                `json:"name"`
	File *multipart.FileHeader `json:"file"`
}
