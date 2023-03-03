package models

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID            string       `json:"aid,omitempty"`    // 'trt13'
	Name          string       `json:"name,omitempty"`   // 'Tribunal Regional do Trabalho 13° Região'
	Type          string       `json:"type,omitempty"`   // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity        string       `json:"entity,omitempty"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF            string       `json:"uf,omitempty"`     // Short code for federative unity.
	URL           string       `json:"url,omitempty"`    // Link for state url
	Collecting    []Collecting `json:"collecting,omitempty"`
	TwitterHandle string       `json:"twitter_handle,omitempty"` // Agency's twitter handle
	OmbudsmanURL  string       `json:"ombudsman_url,omitempty"`  //Agencys's ombudsman url
}

// Collecting A Struct containing the day we checked the status of the data and the reasons why we didn't collected it.
type Collecting struct {
	Timestamp   *int64   `json:"timestamp,omitempty"`   // Day(unix) we checked the status of the data
	Description []string `json:"description,omitempty"` // Reasons why we didn't collect the data
	Collecting  bool     `json:"collecting,omitempty"`  // If there is data from that agency
}
