package backup

// A LogLine represent a json line of LogDNA backup file.
type LogLine struct {
	Container string   `json:"container"`
	Line      string   `json:"_line"`
	Host      string   `json:"_host"`
	Tags      []string `json:"_tag"`
}
