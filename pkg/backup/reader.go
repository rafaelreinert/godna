package backup

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"

	"github.com/apex/log"
)

type Status struct {
	FileInfo os.FileInfo
	Offset   int64
}

type Reader struct {
	File   *os.File
	Status *Status
}

func NewReader(file *os.File) *Reader {
	fs, err := file.Stat()
	if err != nil {
		log.WithError(err).Fatal("Error while geting the stat of input file")
		os.Exit(1)
	}
	s := &Status{FileInfo: fs, Offset: 0}
	return &Reader{File: file, Status: s}
}

func (r *Reader) ReadAll(output chan *LogLine) {
	gzr, err := gzip.NewReader(r.File)
	if err != nil {
		log.WithError(err).Fatal("Error while opening the GZIP reader")
		os.Exit(1)
	}
	defer gzr.Close() // nolint: errcheck
	reader := bufio.NewReaderSize(gzr, 1024*56)
	for {
		var part []byte
		if part, err = reader.ReadSlice('\n'); err != nil {
			if err == io.EOF {
				break
			}
			log.WithError(err).Error("Error while reading the input file")
		}
		var l LogLine
		err := json.Unmarshal(part, &l)
		if err != nil {
			log.WithError(err).Error("Error while unmarshal a line")
		}
		output <- &l
		ps, er := r.File.Seek(0, 1)
		if er != nil {
			log.WithError(er).Error("Error while seeking the input file")
		}
		r.Status.Offset = ps
	}
	close(output)
}
