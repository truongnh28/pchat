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
	Height           uint32 `json:"height,omitempty"`
	Width            uint32 `json:"width,omitempty"`
	FileSize         uint32 `json:"file_size,omitempty"`
	ResourceType     string `json:"resource_type,omitempty"`
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

func (f *File) GetResourceType() string {
	if f == nil {
		return ""
	}
	return f.ResourceType
}

func (f *File) GetFileSize() uint32 {
	if f == nil {
		return 0
	}
	return f.FileSize
}

func (f *File) GetWidth() uint32 {
	if f == nil {
		return 0
	}
	return f.Width
}

func (f *File) GetHeight() uint32 {
	if f == nil {
		return 0
	}
	return f.Height
}
