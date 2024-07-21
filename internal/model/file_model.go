package model

import "mime/multipart"

type File struct {
	Name       string
	FileHeader *multipart.FileHeader `validate:"omitempty,image=500x500"`
}
