# GoDNA [![Build Status](https://travis-ci.org/rafaelreinert/godna.svg?branch=master)](https://travis-ci.org/rafaelreinert/godna) [![codecov](https://codecov.io/gh/rafaelreinert/godna/branch/master/graph/badge.svg)](https://codecov.io/gh/rafaelreinert/godna) [![Go Report Card](https://goreportcard.com/badge/github.com/rafaelreinert/godna)](https://goreportcard.com/report/github.com/rafaelreinert/godna)
GoDNA ia a tool to handle LogDNA backup file.

## Split 

Split the backup in many files per host

```
usage: godna split <file> <output_dir>

split a log.

Flags:
      --help           Show context-sensitive help (also try --help-long and
                       --help-man).
  -c, --containers=CONTAINERS ...  
                       Filter by container.
  -t, --host=HOST ...  Filter by Tag.

Args:
  <file>        File to Split.
  <output_dir>  Folder where GoDNA will save the splited logs.

```

## Exemples

Filter by containers (will save only container matched files):
``` shell
godna split -c postgres -c node ~/Backup/e56cd18d89.2018-04-29.json.gz ~/logs
```

Filter by tags (will save only tags matched files):
``` shell
godna split -t kubernetes -t web ~/Backup/e56cd18d89.2018-04-29.json.gz ~/logs
```

Filter by tags or containers (will save only tags or containers matched files):
``` shell
godna split -t kubernetes -c postgres ~/Backup/e56cd18d89.2018-04-29.json.gz ~/logs
```


