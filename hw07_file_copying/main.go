package main

import (
	"flag"
	"log"
)

var (
	from, to, placeholder string
	limit, offset         int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.StringVar(&placeholder, "placeholder", ".", "placeholder for progressbar")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	switch {
	case limit < 0:
		limit = 0
	case offset < 0:
		offset = 0
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		log.Fatal(err)
	}
}
