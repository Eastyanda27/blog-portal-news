package entity

import (
	"mime/multipart"
)

type FileUploadEntity struct {
	Name       string
	Path       string
	File       multipart.File
	FileHeader *multipart.FileHeader
	Folder     string
}
