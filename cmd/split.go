package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/gosuri/uilive"
	"github.com/rafaelreinert/godna/pkg/backup"
	"github.com/rafaelreinert/godna/pkg/filter"
)

type Split struct {
	File       *os.File
	OutputDir  string
	Containers []string
	Tags       []string
	r          *backup.Reader
	w          *backup.Writer
	filter     *filter.Filter
	ch         chan *backup.LogLine
}

func NewSplit(file *os.File, outputDir string, containers []string, tags []string) *Split {
	fmt.Println("Starting..")
	r := backup.NewReader(file)
	w := backup.NewWriter(outputDir)
	filter := &filter.Filter{Containers: containers, Tags: tags}
	ch := make(chan *backup.LogLine, 1000)
	return &Split{File: file, OutputDir: outputDir, Containers: containers, Tags: tags, r: r, w: w, ch: ch, filter: filter}
}

func (s *Split) Do() {
	fmt.Println("Started")
	var wg sync.WaitGroup
	wg.Add(1)
	go s.r.ReadAll(s.ch)
	go s.followFileReading(&wg)
	for l := range s.ch {
		if s.filter.Do(l) {
			s.w.WriteInFileByServer(l)
		}
	}
	s.w.CloseFiles()
	wg.Wait()
	fmt.Println("Finished")
}

func (s *Split) followFileReading(wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()
	for {
		time.Sleep(100 * time.Millisecond)
		st := s.r.Status
		t := time.Now()
		elapsed := t.Sub(start)
		tf := ((elapsed.Seconds()) * float64(st.FileInfo.Size()-st.Offset)) / float64(st.Offset)
		_, err := fmt.Fprintf(writer, "Processing.. %.2f%% of %s,  estimated time: \t%-v \n", float64(st.Offset*100)/float64(st.FileInfo.Size()), st.FileInfo.Name(), time.Duration(tf)*time.Second)
		if err != nil {
			log.WithError(err).Error("Error while updating the status")
		}
		if st.Offset == st.FileInfo.Size() {
			break
		}
	}

}
