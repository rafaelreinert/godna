package backup

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteInFileByServer(t *testing.T) {
	dir, err := ioutil.TempDir("", "logdna-logs")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	log1 := LogLine{Host: "abc-123", Line: "hi"}
	log2 := LogLine{Host: "zyw-987", Line: "bye"}
	log3 := LogLine{Host: "abc-123", Line: "hello"}
	w := NewWriter(dir)
	w.WriteInFileByServer(&log1)
	w.WriteInFileByServer(&log2)
	w.WriteInFileByServer(&log3)
	w.CloseFiles()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) != 2 {
		t.Errorf("expected only 2 file but found %d.", len(files))
	}
	var contents []string
	for _, file := range files {
		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		contents = append(contents, string(content))
	}
	if files[0].Name() != "abc-123.log" && contents[0] != "hi\nhello" {
		t.Errorf("%v was unexpected.", files[0].Name())
	}
	if contents[0] != "hi\nhello\n" {
		t.Errorf("%v was wrong %v.", files[0].Name(), contents[0])
	}
	if files[1].Name() != "zyw-987.log" {
		t.Errorf("%v was unexpected.", files[1].Name())
	}
	if contents[1] != "bye\n" {
		t.Errorf("%v was wrong %v.", files[1].Name(), contents[1])
	}
}
