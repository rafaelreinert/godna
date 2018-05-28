package main

import (
	"os"

	"github.com/rafaelreinert/godna/cmd"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	var app = kingpin.New("godna", "A logdna backup files mananger application.")
	var split = app.Command("split", "split a log.")
	var file = split.Arg("file", "File to Split.").Required().File()
	var outputDir = split.Arg("output_dir", "Folder where GoDNA will save the splited logs.").Required().ExistingDir()
	var containers = app.Flag("containers", "Filter by container.").Short('c').Strings()
	var tags = app.Flag("host", "Filter by Tag.").Short('t').Strings()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case split.FullCommand():
		s := cmd.NewSplit(*file, *outputDir, *containers, *tags)
		s.Do()
	}
}
