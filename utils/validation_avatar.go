package utils

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func ValidateAvatar(fileheader *multipart.FileHeader) (string, error) {
	if (fileheader.Size) > 100*1024 {
		return "", errors.New("file size is too large, max 100 kb is alloweds")
	}

	fileExtension := strings.ToLower(filepath.Ext(fileheader.Filename))
	if fileExtension != ".jpg" && fileExtension != ".jpeg" {
		return fileExtension, errors.New("invalid file type")
	}

	return fileExtension, nil
}
