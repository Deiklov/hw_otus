package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// прочитать со смещением
	// сравнить его с размером файла если что выходим ErrOffsetExceedsFileSize
	// если нет длины то ErrUnsupportedFile
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The file does not exist.")
		} else {
			fmt.Println("Error:", err)
		}
		return err
	}

	// Check if the file is a regular file
	if !fileInfo.Mode().IsRegular() {
		fmt.Println("Cannot get the length of a non-regular file.")
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	input, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer output.Close()

	// Create a new progress bar with the source file size as the total
	bar := pb.Full.Start64(fileInfo.Size() - offset)

	_, err = input.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Create a proxy reader that reports progress to the bar
	barReader := bar.NewProxyReader(input)

	defer bar.Finish()

	if limit == 0 {
		_, err = io.Copy(output, barReader)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	} else {
		_, err = io.CopyN(output, barReader, limit)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}
