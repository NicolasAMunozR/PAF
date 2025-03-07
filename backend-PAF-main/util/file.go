package util

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
)

// SaveCompressed : Se comprime imagen
func SaveCompressed(src image.Image, quality int, absolutePath string) error {
	f, err := os.Create(absolutePath)
	log.Printf("Creado %v", f)
	log.Printf("err %v", err)

	if err != nil {
		return err
	}
	defer f.Close()
	err = imaging.Encode(f, src, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		return err
	}
	f.Sync()
	log.Print("Sincronizado")
	return nil
}

func UnZipAndSaveFile(file *multipart.FileHeader, dir string) (string, string, error) {

	var err error
	newFileName := uuid.Must(uuid.NewV4(), err)
	if err != nil {
		return "", "", err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}
	fileData, err := file.Open()
	fileExtension := filepath.Ext(file.Filename)

	absolutePath := fmt.Sprintf("%s/%s%s", dir, newFileName.String(), fileExtension)

	if err != nil {
		log.Printf("error %v", err)
		return "", "", err
	}
	fileBytes, err := ioutil.ReadAll(fileData)
	if err != nil {
		log.Printf("error %v", err)
		return "", "", err
	}

	var buf *bytes.Buffer = bytes.NewBuffer(fileBytes)
	gRead, err := zlib.NewReader(buf)

	if err != nil {
		log.Printf("error %v", err)
		return "", "", err
	}
	gRead.Close()
	writer, err := os.Create(absolutePath)
	if err != nil {
		log.Printf("error %v", err)
		return "", "", err
	}
	_, err = io.Copy(writer, gRead)
	if err != nil {
		log.Printf("error %v", err)
		return "", "", err
	}
	writer.Close()
	var fileType string
	if fileExtension == ".mp4" {
		fileType = "video"
	} else if fileExtension == ".pdf" {
		fileType = "pdf"
	} else {
		fileType = "image"
	}
	return absolutePath, fileType, nil
}

func SaveVideo(file *multipart.FileHeader, dir string) (string, error) {

	var err error

	newFileName := uuid.Must(uuid.NewV4(), err)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}

	fileExtension := filepath.Ext(file.Filename)

	relativePath := fmt.Sprintf("%s%s", newFileName.String(), fileExtension)
	absolutePath := fmt.Sprintf("%s/%s%s", dir, newFileName.String(), fileExtension)
	fileData, err := file.Open()
	if err != nil {
		log.Printf("error %v", err)
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(fileData)
	if err != nil {
		log.Printf("error %v", err)
		return "", err
	}
	ioutil.WriteFile(absolutePath, fileBytes, 0777)

	return relativePath, nil
}

func SaveFile(file *multipart.FileHeader, dir string) (string, string, error) {
	var err error

	newFileName := uuid.Must(uuid.NewV4(), err)
	if err != nil {
		return "", "", err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}
	fileExtension := filepath.Ext(file.Filename)
	resizedRelativePath := fmt.Sprintf("resized-%s%s", newFileName.String(), fileExtension)
	resizedPath := fmt.Sprintf("%s/resized-%s%s", dir, newFileName.String(), fileExtension)

	relativePath := fmt.Sprintf("%s%s", newFileName.String(), fileExtension)
	absolutePath := fmt.Sprintf("%s/%s%s", dir, newFileName.String(), fileExtension)

	fileData, err := file.Open()

	if err != nil {

		return "", "", err
	}
	src, err := imaging.Decode(fileData, imaging.AutoOrientation(true))

	if err != nil {
		return "", "", err
	}

	err = SaveCompressed(src, 55, absolutePath)

	if err != nil {
		return "", "", err
	}
	src, err = imaging.Open(absolutePath, imaging.AutoOrientation(true))
	if err != nil {
		return "", "", err
	}
	src = imaging.Fill(src, 150, 150, imaging.Center, imaging.Lanczos)

	err = imaging.Save(src, resizedPath)

	if err != nil {
		return "", "", err
	}
	return relativePath, resizedRelativePath, err
}
