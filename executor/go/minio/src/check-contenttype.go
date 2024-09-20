package src

import "path/filepath"

func GetContentType(objName string) string {
	var contentType string

	fileExtension := filepath.Ext(objName)
	if fileExtension == ".zip" {
		contentType = "application/zip"
	} else if fileExtension == ".pdf" {
		contentType = "application/pdf"
	}

	return contentType
}
