package filter

import (
	"testing"

	"github.com/rafaelreinert/godna/pkg/backup"
)

// nolint: gocyclo
func TestDoWithEmptyFilters(t *testing.T) {
	containers := []string{}
	tags := []string{}
	filter := Filter{Containers: containers, Tags: tags}

	var lines []backup.LogLine
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-01", Tags: []string{"app"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-02", Tags: []string{"app"}})

	for _, v := range lines {
		if !filter.Do(&v) {
			t.Errorf("%v was filtered.", v)
		}
	}

}

// nolint: gocyclo
func TestDoWithContainerFilters(t *testing.T) {
	containers := []string{"mysql", "grafana"}
	tags := []string{}
	filter := Filter{Containers: containers, Tags: tags}

	var lines []backup.LogLine
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-01", Tags: []string{"app"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-02", Tags: []string{"app"}})

	for _, v := range lines {
		if filter.Do(&v) {
			if v.Container != "mysql" && v.Container != "grafana" {
				t.Errorf("%v was not filtered.", v)
			}
		}
	}

}

// nolint: gocyclo
func TestDoWithTagsFilters(t *testing.T) {
	containers := []string{}
	tags := []string{"db"}
	filter := Filter{Containers: containers, Tags: tags}

	var lines []backup.LogLine
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-01", Tags: []string{"app"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-02", Tags: []string{"app"}})

	for _, v := range lines {
		if filter.Do(&v) {
			if v.Container != "mysql" && v.Container != "postgres" {
				t.Errorf("%v was not filtered.", v)
			}
		}
	}

}

func TestDoWithContainersAndTagsFilters(t *testing.T) {
	containers := []string{"mysql"}
	tags := []string{"app"}
	filter := Filter{Containers: containers, Tags: tags}

	var lines []backup.LogLine
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "mysql", Line: "Init mysql", Host: "mysql-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-01", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "postgres", Line: "Init postgres", Host: "postgres-02", Tags: []string{"db"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-01", Tags: []string{"app"}})
	lines = append(lines, backup.LogLine{Container: "grafana", Line: "Init grafana", Host: "grafana-02", Tags: []string{"app"}})

	for _, v := range lines {
		if filter.Do(&v) {
			if v.Container != "mysql" && v.Container != "grafana" {
				t.Errorf("%v was not filtered.", v)
			}
		}
	}

}
