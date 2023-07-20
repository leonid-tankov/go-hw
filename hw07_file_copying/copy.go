package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeArgument      = errors.New("one of the arguments is negative")
	BufferSize               = 100
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return ErrNegativeArgument
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	fileSize := fileInfo.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	src, err := os.OpenFile(fromPath, os.O_RDONLY, 0o444)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	if (limit > fileSize-offset) || (limit == 0) {
		limit = fileSize - offset
	}
	bSize := BufferSize
	if BufferSize > int(limit) {
		bSize = int(limit)
	}
	buf := make([]byte, bSize)
	fmt.Println("Start copying")
	for i := 1; i <= int(limit); i += bSize {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dst.Write(buf[:n]); err != nil {
			return err
		}
		progressBytes := i
		percent := 100 * float64(i) / float64(limit)
		if i > int(limit)-bSize {
			percent = 100
			progressBytes = int(limit)
		}
		progress := strings.Repeat("#", int(percent)/2)
		fmt.Printf("\r[%-50s]%5.1f%% %15d/%d bytes ", progress, percent, progressBytes, limit)
	}
	fmt.Println("\nComplete copying")
	return nil
}
