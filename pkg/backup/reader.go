package backup

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"

	"github.com/apex/log"
)

// A Status contain  the status of reading.
// FileInfo is a os.FileInfo of reading file.
// Offset is a position where the reader read at moment
type Status struct {
	FileInfo os.FileInfo
	Offset   int64
}

// A Reader can be read to retrieve
// uncompressed data from logdna backup file compressed with gzip.
type Reader struct {
	File   *os.File
	Status *Status
}

// NewReader creates a new Reader reading the given file.
func NewReader(file *os.File) *Reader {
	fs, err := file.Stat()
	if err != nil {
		log.WithError(err).Fatal("Error while geting the stat of input file")
		os.Exit(1)
	}
	s := &Status{FileInfo: fs, Offset: 0}
	return &Reader{File: file, Status: s}
}

// ReadAll reads all file until EOF,
// Each row is parsed to a LogLine and sent to output channel.
// After EOF the output channel will be closed
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
			if err == io.EOF && len(part) == 0 {
				break
			} else if err != io.EOF {
				log.WithError(err).Error("Error while reading the input file")
			}
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
