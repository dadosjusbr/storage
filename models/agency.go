package models

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID            string       `json:"aid"`    // 'trt13'
	Name          string       `json:"name"`   // 'Tribunal Regional do Trabalho 13° Região'
	Type          string       `json:"type"`   // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity        string       `json:"entity"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF            string       `json:"uf"`     // Short code for federative unity.
	URL           string       `json:"url"`    // Link for state url
	Collecting    []Collecting `json:"collecting"`
	TwitterHandle string       `json:"twitter_handle"` // Agency's twitter handle
	Ombudsman     string       `json:"ombudsman"`      //Agencys's ombudsman
}

// Collecting A Struct containing the day we checked the status of the data and the reasons why we didn't collected it.
type Collecting struct {
	Timestamp   *int64   `json:"timestamp"`   // Day(unix) we checked the status of the data
	Description []string `json:"description"` // Reasons why we didn't collect the data
}
