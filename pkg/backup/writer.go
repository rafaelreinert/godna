package backup

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
)

// A Writer can be write to text plain from LogLine.
type Writer struct {
	OutputDir string
	files     map[string]*os.File
}

// NewWriter creates a new Writer writing files in output dir.
func NewWriter(outputDir string) *Writer {
	return &Writer{OutputDir: outputDir, files: make(map[string]*os.File)}
}

// WriteInFileByServer write a LogLine in a respective file.
// The file name is {l.Host}.log.
// For Each Distinct Host will open a new file but will not close the files when done.
// It is the caller's responsibility to call CloseFiles on the Writer when done.
func (w *Writer) WriteInFileByServer(l *LogLine) {
	var f *os.File
	var err error
	f, ok := w.files[l.Host]
	if !ok {
		f, err = os.Create(filepath.Join(w.OutputDir, l.Host+".log"))
		if err != nil {
			log.WithError(err).Fatal("Error while creating a output file")
			os.Exit(1)
		}
		w.files[l.Host] = f
	}
	_, err = f.Write([]byte(l.Line + "\n"))
	if err != nil {
		log.WithError(err).Error("Error while writing the output file")
	}

}

// CloseFiles will close all opened files.
func (w *Writer) CloseFiles() {
	for _, f := range w.files {
		err := f.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing a output file")
		}
	}
}
