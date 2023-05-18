package main

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/tonbyte/provider-node/datatype"
)

func checkFileSize(file *multipart.FileHeader, freeSpaceLeftInKb int64) error {
	if file.Size == 0 || file.Size > freeSpaceLeftInKb*1024 {
		return errors.New("not enough space")
	}

	return nil
}

func saveFile(filesPath string, freeSpaceLeft int64, files []*multipart.FileHeader) (datatype.FileInfo, error) {
	fileInfo := datatype.FileInfo{}

	file := files[0]
	err := checkFileSize(file, freeSpaceLeft)
	if err != nil {
		return datatype.FileInfo{}, err
	}

	var src multipart.File
	src, err = file.Open()
	if err != nil {
		return datatype.FileInfo{}, err
	}
	defer src.Close()

	var dst *os.File
	fileInfo.FullPath = filepath.Join(filesPath, file.Filename)
	dst, err = os.Create(fileInfo.FullPath)
	if err != nil {
		return datatype.FileInfo{}, err
	}
	defer dst.Close()

	fileInfo.Size = strconv.FormatInt(file.Size, 10)
	if _, err = io.Copy(dst, src); err != nil {
		return datatype.FileInfo{}, err
	}

	return fileInfo, nil
}
