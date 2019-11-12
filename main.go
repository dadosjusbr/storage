package main

import (
	"fmt"
	"log"

	d "github.com/dadosjusbr/storage/db"
)

func main() {
	fmt.Println("Starting the application...")
	mppb := d.Agency{ID: nil, ShortName: "HH", Name: "Tribunal Eleitoral da Paraíba", Type: "E", Entity: "J", UF: "PB"}
	C, err := d.NewClient("mongodb://localhost:27017", "test")
	if err != nil {
		log.Fatalf("Error trying to connect to db: %q", err)
	}
	C.SaveAgency(mppb)
}
