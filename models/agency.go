package models

// Agency A Struct containing the main descriptions of each Agency.
type Agency struct {
	ID         string       `json:"id" bson:"aid,omitempty"`          // 'trt13'
	Name       string       `json:"nome" bson:"name,omitempty"`       // 'Tribunal Regional do Trabalho 13° Região'
	Type       string       `json:"jurisdicao" bson:"type,omitempty"` // "R" for Regional, "M" for Municipal, "F" for Federal, "E" for State.
	Entity     string       `json:"entidade" bson:"entity,omitempty"` // "J" For Judiciário, "M" for Ministério Público, "P" for Procuradorias and "D" for Defensorias.
	UF         string       `json:"uf" bson:"uf,omitempty"`           // Short code for federative unity.
	FlagURL    string       `json:"url" bson:"url,omitempty"`         // Link for state url
	Collecting []Collecting `json:"coletando" bson:"collecting,omitempty"`
}

// Collecting A Struct containing the day we checked the status of the data and the reasons why we didn't collected it.
type Collecting struct {
	Timestamp   *int64   `json:"timestamp" bson:"timestamp,omitempty"`   // Day(unix) we checked the status of the data
	Description []string `json:"descricao" bson:"description,omitempty"` // Reasons why we didn't collect the data
}
