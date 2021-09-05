package domain

import "io"

type ImageFile struct {
	FileName string
	Content io.Reader
}

func NewImageFile(fn string, buf io.Reader) *ImageFile{
	return &ImageFile{
		FileName:  fn,
		Content: buf,
	}
}