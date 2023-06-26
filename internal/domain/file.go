package domain

import "mime/multipart"

type UploadIn struct {
	FileName string         `json:"file_name,omitempty"`
	FileData multipart.File `json:"file_data,omitempty"`
}

type File struct {
	Id               uint   `json:"id,omitempty"`
	SecureURL        string `json:"secure_url"`
	OriginalFilename string `json:"original_filename"`
}

func (f *File) GetSecureURL() string {
	if f == nil {
		return ""
	}
	return f.SecureURL
}

func (f *File) GetOriginalFileName() string {
	if f == nil {
		return ""
	}
	return f.OriginalFilename
}

func (f *File) GetId() uint {
	if f == nil {
		return 0
	}
	return f.Id
}
