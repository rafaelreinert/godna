package backup

type LogLine struct {
	Container string   `json:"container"`
	Line      string   `json:"_line"`
	Host      string   `json:"_host"`
	Tags      []string `json:"_tag"`
}
