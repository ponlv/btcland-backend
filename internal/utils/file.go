package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/h2non/filetype"
)

func GetFileExtension(fileHeader *multipart.FileHeader) (string, error) {
	// Mở tệp
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Use filetype library to detect the MIME type
	kind, err := filetype.Match(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to detect file type: %w", err)
	}

	if kind == filetype.Unknown {
		return "", fmt.Errorf("không tìm thấy phần mở rộng cho MIME type")
	}

	return fmt.Sprintf(".%s", kind.Extension), nil
}
