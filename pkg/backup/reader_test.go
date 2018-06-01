package backup

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
	l2 = "{\"container\":\"appy\",\"_host\":\"appy-fn23p\",\"_line\":\"append x\",\"_tag\":[\"web\"]}"
)

func createTempFile() (string, string) {
	var buf bytes.Buffer
	dir, err := ioutil.TempDir("", "logdna-backup")
	if err != nil {
		log.Fatal(err)
	}
	content := []byte(l1)
	content = append(content, []byte(l2)...)

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

// nolint: gocyclo
func TestReadAll(t *testing.T) {
	name, dir := createTempFile()
	defer os.RemoveAll(dir)
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := NewReader(f)
	c := make(chan *LogLine)
	go r.ReadAll(c)
	var lines []*LogLine
	for l := range c {
		lines = append(lines, l)
	}
	l1 := *lines[0]
	if l1.Container != "appx" || l1.Host != "appx-wx6pv" || l1.Line != "init download" || l1.Tags[0] != "job" {
		t.Errorf("%v was unexpected in first line.", l1)
	}

	l2 := *lines[1]
	if l2.Container != "appy" || l2.Host != "appy-fn23p" || l2.Line != "append x" || l2.Tags[0] != "web" {
		t.Errorf("%v was unexpected in second line.", l2)
	}

	if len(lines) != 2 {
		t.Errorf("expected only 2 lines but found %d.", len(lines))
	}

}
