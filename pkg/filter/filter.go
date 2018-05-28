package filter

import "github.com/rafaelreinert/godna/pkg/backup"

type Filter struct {
	Containers []string
	Tags       []string
}

func (f *Filter) Do(line *backup.LogLine) bool {
	if len(f.Containers) > 0 || len(f.Tags) > 0 {
		if len(f.Containers) > 0 && include(f.Containers, line.Container) {
			return true
		}
		if len(f.Tags) > 0 && any(f.Tags, line.Tags) {
			return true
		}
		return false
	}
	return true
}

func include(vs []string, t string) bool {
	for _, v := range vs {
		if v == t {
			return true
		}
	}
	return false
}

func any(vs []string, svs []string) bool {
	for _, v := range vs {
		if include(svs, v) {
			return true
		}
	}
	return false
}
