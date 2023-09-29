package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// создаем временный файл
	tmpFile, err := os.CreateTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}

	// не забываем закрыть и удалить временные файлы
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(tmpFile.Name())
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(tmpFile)

	t.Run("empty file", func(t *testing.T) {
		fromFile := "testdata/empty.txt"
		err := Copy(fromFile, tmpFile.Name(), 0, 0)
		require.Nil(t, err)
		tmpFileStat, err := tmpFile.Stat()
		require.Nil(t, err)
		require.Equal(t, int64(0), tmpFileStat.Size())
	})

	t.Run("invalid input file", func(t *testing.T) {
		fromFile := "testdata/invalid.txt"
		err := Copy(fromFile, tmpFile.Name(), 6000, 0)
		require.NotNil(t, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("invalid output file", func(t *testing.T) {
		fromFile := "testdata/input.txt"
		outFile := "/root/outfile.txt"
		err := Copy(fromFile, outFile, 0, 0)
		require.NotNil(t, err)
		require.Equal(t, ErrOpenWriteFile, err)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		fromFile := "testdata/input.txt"
		err := Copy(fromFile, tmpFile.Name(), 7000, 0)
		require.NotNil(t, err)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})
}
