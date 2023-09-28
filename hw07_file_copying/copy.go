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
	size, err := validate(fromPath, toPath, offset, &limit)
	if err != nil {
		return err
	}
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
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
		if closeErr := dst.Close(); closeErr != nil {
			log.Fatal("failed close file")
		}
	}()

	return processCopy(src, dst, offset, limit, size)
}

func validate(fromPath, toPath string, offset int64, limit *int64) (int64, error) {
	if from == "" || toPath == "" {
		return 0, ErrUnsupportedFile
	}

	if *limit != 0 {
		*limit += offset
	}

	// статистика по файлу источнику
	fromFileStat, err := os.Stat(fromPath)
	if err != nil {
		return 0, err
	}

	// если исходный файл не обычный - кидаем ошибку
	if !fromFileStat.Mode().IsRegular() {
		return 0, ErrIsNotRegularFile
	}

	// если сдвиг больше размера исходного файла - кидаем ошибку
	size := fromFileStat.Size()
	if size < offset {
		return size, ErrOffsetExceedsFileSize
	}

	return size, nil
}

func processCopy(src, dst *os.File, offset, limit, fromFileSize int64) error {
	// создаем прогрессбар
	bar := Bar{}
	pbSize := fromFileSize
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
