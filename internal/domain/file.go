package domain

import "mime/multipart"

type UploadIn struct {
	FileName string         `json:"file_name,omitempty"`
	FileData multipart.File `json:"file_data,omitempty"`
}
