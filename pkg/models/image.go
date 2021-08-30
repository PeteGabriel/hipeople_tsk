package models

import "mime/multipart"

type ImageFile struct {
	Header *multipart.FileHeader
	Content multipart.File
}

func NewImageFile(h *multipart.FileHeader, c multipart.File) *ImageFile{
	return &ImageFile{
		Header:  h,
		Content: c,
	}
}