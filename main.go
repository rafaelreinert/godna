package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/gosuri/uilive"
	"gopkg.in/alecthomas/kingpin.v2"
)

type logLine struct {
	Container string   `json:"container"`
	Line      string   `json:"_line"`
	Host      string   `json:"_host"`
	Tags      []string `json:"_tag"`
}

var (
	app        = kingpin.New("godna", "A logdna backup files mananger application.")
	split      = app.Command("split", "split a log.")
	file       = split.Arg("file", "File to Split.").Required().File()
	outputDir  = split.Arg("output_dir", "Folder where GoDNA will save the splited logs.").Required().ExistingDir()
	containers = app.Flag("containers", "Filter by container.").Short('c').Strings()
	tags       = app.Flag("host", "Filter by Tag.").Short('t').Strings()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case split.FullCommand():
		splitFile(*file, *outputDir, *containers, *tags)
	}
}

func splitFile(file *os.File, outputDir string, containers []string, tags []string) {
	fmt.Println("Starting..")
	ch := make(chan *logLine, 1000)
	filteredCh := make(chan *logLine, 1000)
	go readFile(ch, file)
	go filterLines(ch, filteredCh, containers, tags)
	fmt.Println("Started")
	writeFile(filteredCh, outputDir)
	fmt.Println("Finished")
}

func filterLines(ch chan *logLine, fch chan *logLine, containers []string, tags []string) {
	for log := range ch {
		if len(containers) > 0 || len(tags) > 0 {
			if len(containers) > 0 && include(containers, log.Container) {
				fch <- log
			}
			if len(tags) > 0 && any(tags, log.Tags) {
				fch <- log
			}
		} else {
			fch <- log
		}
	}
	close(fch)
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

func any(vs []string, svs []string) bool {
	for _, v := range vs {
		if include(svs, v) {
			return true
		}
	}
	return false
}

func followFileReading(ch chan int, f *os.File) {
	start := time.Now()
	fs, err := f.Stat()
	if err != nil {
		log.WithError(err).Fatal("Error while geting the stat of input file")
		os.Exit(1)
	}
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	for _ = range ch {
		ps, err := f.Seek(0, 1)
		t := time.Now()
		elapsed := t.Sub(start)
		tf := ((elapsed.Seconds()) * float64(fs.Size()-ps)) / float64(ps)
		if err != nil {
			log.WithError(err).Error("Error while seeking the input file")
		}
		fmt.Fprintf(writer, "Processing.. %.2f%% of %s, estimated time: \t%v\n", float64(ps*100)/float64(fs.Size()), fs.Name(), time.Duration(tf)*time.Second)
	}
	fmt.Fprintf(writer, "Finished: Processed 100%% of %s\n", fs.Name())
}

func readFile(ch chan *logLine, f *os.File) {
	r, err := gzip.NewReader(f)
	if err != nil {
		log.WithError(err).Fatal("Error while opening the GZIP reader")
		os.Exit(1)
	}
	defer r.Close()
	reader := bufio.NewReaderSize(r, 1024*56)
	frch := make(chan int, 100)
	go followFileReading(frch, f)
	i := 0
	for {
		var part []byte
		if part, err = reader.ReadSlice('\n'); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			log.WithError(err).Error("Error while reading the input file")
		}
		var l logLine
		err := json.Unmarshal(part, &l)
		if err != nil {
			log.WithError(err).Error("Error while unmarshal a line")
		}
		ch <- &l
		i++
		if i%10000 == 0 {
			frch <- i
		}
	}
	close(ch)
	close(frch)
}

func writeFile(chw chan *logLine, outputDir string) {
	fs := make(map[string]*os.File)

	for l := range chw {
		var f *os.File
		var err error
		f, ok := fs[l.Host]
		if !ok {
			f, err = os.Create(outputDir + "/" + l.Host + ".log")
			if err != nil {
				log.WithError(err).Fatal("Error while creating a output file")
				os.Exit(1)
			}
			fs[l.Host] = f
		}
		_, err = f.Write([]byte(l.Line + "\n"))
		if err != nil {
			log.WithError(err).Error("Error while writing the output file")
			return
		}
	}

	defer closeFiles(fs)
}

func closeFiles(fs map[string]*os.File) {
	for _, f := range fs {
		f.Close()
	}
}
