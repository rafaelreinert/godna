package cmd

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	l1 = "{\"container\":\"appx\",\"_host\":\"appx-wx6pv\",\"_line\":\"init download\",\"_tag\":[\"job\"]}\n"
	l2 = "{\"container\":\"appy\",\"_host\":\"appy-fn23p\",\"_line\":\"append x\",\"_tag\":[\"web\"]}\n"
	l3 = "{\"container\":\"appz\",\"_host\":\"appy-za2wq\",\"_line\":\"append y\",\"_tag\":[\"mobile\"]}"
)

func createTempFile() (string, string) {
	var buf bytes.Buffer
	dir, err := ioutil.TempDir("", "logdna-backup")
	if err != nil {
		log.Fatal(err)
	}
	content := []byte(l1)
	content = append(content, []byte(l2)...)
	content = append(content, []byte(l3)...)

	zw := gzip.NewWriter(&buf)

	_, err = zw.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	tmpfn := filepath.Join(dir, "logdna.json.gz")
	if err := ioutil.WriteFile(tmpfn, buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
	return tmpfn, dir
}

func createOutputDir() string {
	dir, err := ioutil.TempDir("", "logdna-logs")
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// nolint: gocyclo
func TestSplit(t *testing.T) {
	dir := createOutputDir()
	defer os.RemoveAll(dir)
	name, inDir := createTempFile()
	defer os.RemoveAll(inDir)
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := NewSplit(f, dir, []string{"appx"}, []string{"web"})
	s.Do()
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
	if files[0].Name() != "appx-wx6pv.log" && contents[0] != "hi\nhello" {
		t.Errorf("%v was unexpected.", files[0].Name())
	}
	if contents[0] != "init download\n" {
		t.Errorf("%v was wrong %v.", files[0].Name(), contents[0])
	}
	if files[1].Name() != "appy-fn23p.log" {
		t.Errorf("%v was unexpected.", files[1].Name())
	}
	if contents[1] != "append x\n" {
		t.Errorf("%v was wrong %v.", files[1].Name(), contents[1])
	}

}
