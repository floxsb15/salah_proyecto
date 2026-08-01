package controllers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	maxUserPhotoSize    = 2 << 20
	maxVehiclePhotoSize = 5 << 20
	maxGuaranteeSize    = 10 << 20
)

var imageMIMEExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

var guaranteeMIMEExtensions = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
}

func validateImageUpload(header *multipart.FileHeader, maxSize int64) (string, error) {
	return validateUpload(header, maxSize, imageMIMEExtensions)
}

func validateGuaranteeUpload(header *multipart.FileHeader) (string, error) {
	return validateUpload(header, maxGuaranteeSize, guaranteeMIMEExtensions)
}

func validateUpload(header *multipart.FileHeader, maxSize int64, allowed map[string]string) (string, error) {
	if header == nil || header.Size <= 0 || header.Size > maxSize {
		return "", errors.New("Tamano de archivo no permitido")
	}
	file, err := header.Open()
	if err != nil {
		return "", errors.New("No se pudo validar el archivo")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	read, readErr := io.ReadFull(file, buffer)
	if readErr != nil && readErr != io.ErrUnexpectedEOF {
		return "", errors.New("No se pudo validar el archivo")
	}
	mimeType := http.DetectContentType(buffer[:read])
	extension, ok := allowed[mimeType]
	if !ok {
		return "", errors.New("Tipo de archivo no permitido")
	}
	return extension, nil
}

func createUploadFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
}
