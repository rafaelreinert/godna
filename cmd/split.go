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
}

func (s *Split) Do() {
	fmt.Println("Starting..")
	s.r = backup.NewReader(s.File)
	ch := make(chan *backup.LogLine, 1000)
	go s.r.ReadAll(ch)
	fmt.Println("Started")
	var wg sync.WaitGroup
	wg.Add(1)
	go s.followFileReading(&wg)
	filter := filter.Filter{Containers: s.Containers, Tags: s.Tags}
	w := backup.NewWriter(s.OutputDir)
	for l := range ch {
		if filter.Do(l) {
			w.WriteInFileByServer(l)
		}
	}
	w.CloseFiles()
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
		_, err := fmt.Fprintf(writer, "Processing.. %.2f%% of %s, estimated time: \t%v\r\n", float64(st.Offset*100)/float64(st.FileInfo.Size()), st.FileInfo.Name(), time.Duration(tf)*time.Second)
		if err != nil {
			log.WithError(err).Error("Error while updating the status")
		}
		if st.Offset == st.FileInfo.Size() {
			_, err := fmt.Fprintf(writer, "Finished: Processed 100%% of %s\n", st.FileInfo.Name())
			if err != nil {
				log.WithError(err).Error("Error while updating the status")
			}
			break
		}
	}

}
