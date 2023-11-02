package models

// IndexInformation A struct contains a summary of the agency's indexes and their metadata
type IndexInformation struct {
	Month int    `json:"mes,omitempty"`
	Year  int    `json:"ano,omitempty"`
	Meta  *Meta  `json:"meta,omitempty"`
	Score *Score `json:"score,omitempty"`
	Type  string `json:"jurisdicao,omitempty"`
}
