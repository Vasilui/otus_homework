package main

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrIsNotRegularFile      = errors.New("is not regular file")
	BufferSize               = int64(10)
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenReadFile          = errors.New("failed open file for read")
	ErrOpenWriteFile         = errors.New("failed open file for write")
	ErrorEqualFiles          = errors.New("using the same file ")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	size, err := validate(fromPath, toPath, offset, &limit)
	if err != nil {
		return err
	}
	src, err := os.Open(fromPath)
	if err != nil {
		return ErrOpenReadFile
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
	if fromPath == "" || toPath == "" {
		return 0, ErrUnsupportedFile
	}

	if *limit != 0 {
		*limit += offset
	}

	// статистика по файлу источнику
	fromFileStat, err := os.Stat(fromPath)
	if err != nil {
		return 0, ErrUnsupportedFile
	}

	// статистика по файлу на запись
	toFileStat, err := os.Stat(toPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return 0, ErrUnsupportedFile
	}

	if os.SameFile(fromFileStat, toFileStat) {
		return 0, ErrorEqualFiles
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
	if limit != 0 && limit < fromFileSize {
		pbSize = limit
	} else {
		limit = pbSize
	}

	bar.NewOption(offset, pbSize)

	for {
		if offset+BufferSize > limit {
			BufferSize = limit - offset
			offset = limit
		} else {
			offset += BufferSize
		}

		n, errCopy := io.CopyN(dst, src, BufferSize)
		if errCopy != nil && !errors.Is(errCopy, io.EOF) {
			return errCopy
		}

		if n == 0 {
			break
		}

		bar.Play(offset)
	}

	bar.Finish()

	return nil
}
