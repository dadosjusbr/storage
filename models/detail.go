package models

// Detail A struct contains a summary of the agency's indices and their metadata
type Detail struct {
	Month int    `json:"mes,omitempty"`
	Year  int    `json:"ano,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
	Score *Score `json:"score,omitempty"`
}
