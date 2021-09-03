package domain

import "mime/multipart"

type ImageFile struct {
	FileName string
	Content multipart.File
}

func NewImageFile(fn string, c multipart.File) *ImageFile{
	return &ImageFile{
		FileName:  fn,
		Content: c,
	}
}