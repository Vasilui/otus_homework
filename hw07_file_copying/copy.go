package main

import (
	"errors"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrIsNotRegularFile      = errors.New("is not regular file")
	BufferSize               = 10
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenWriteFile         = errors.New("failed open file for write")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if from == "" || to == "" {
		return ErrUnsupportedFile
	}

	if limit != 0 {
		limit += offset
	}

	// статистика по файлу источнику
	fromFileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	// если исходный файл не обычный - кидаем ошибку
	if !fromFileStat.Mode().IsRegular() {
		return ErrIsNotRegularFile
	}

	// если сдвиг больше размера исходного файла - кидаем ошибку
	if fromFileStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := src.Close()
		if closeErr != nil {
			log.Fatal("failed close file")
		}
	}()

	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return ErrOpenWriteFile
	}

	defer func() {
		closeErr := dst.Close()
		if closeErr != nil {
			log.Fatal("failed close file")
		}
	}()

	// создаем прогрессбар
	bar := Bar{}
	pbSize := fromFileStat.Size()
	if limit != 0 {
		pbSize = limit
	}

	bar.NewOption(offset, pbSize)

	buf := make([]byte, BufferSize)
	for {
		n, errRead := src.Read(buf)
		if errRead != nil && errRead != io.EOF {
			return errRead
		}

		if n == 0 {
			break
		}

		if limit != 0 && offset+int64(n) > limit {
			n = int(limit - offset)
			offset = limit
		} else {
			offset += int64(n)
		}

		bar.Play(offset)

		if _, errWrite := dst.Write(buf[:n]); errWrite != nil {
			return errWrite
		}
	}

	bar.Finish()

	return nil
}
