package backup

import (
	"os"

	"github.com/apex/log"
)

type Writer struct {
	OutputDir string
	files     map[string]*os.File
}

func NewWriter(outputDir string) *Writer {
	return &Writer{OutputDir: outputDir, files: make(map[string]*os.File)}
}

func (w *Writer) WriteInFileByServer(l *LogLine) {
	var f *os.File
	var err error
	f, ok := w.files[l.Host]
	if !ok {
		f, err = os.Create(w.OutputDir + "/" + l.Host + ".log")
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

func (w *Writer) CloseFiles() {
	for _, f := range w.files {
		err := f.Close()
		if err != nil {
			log.WithError(err).Error("Error while closing a output file")
		}
	}
}
