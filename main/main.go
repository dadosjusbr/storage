package main

import (
	"encoding/json"
	"fmt"

	storage "github.com/dadosjusbr/storage"
)

func main() {
	c, err := storage.NewDBClient("mongodb+srv://storage:As9t4sDZTqAhfGFV@dadosjusbr-xwain.gcp.mongodb.net/test?retryWrites=true&w=majority", "db", "mi", "agency")
	if err != nil {
		panic(err)
	}
	Client, err := storage.NewClient(c, nil)
	agsMR, _ := Client.GetDataForSecondScreen(9, 2018, "mppb")
	_, err = json.MarshalIndent(agsMR, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(agsMR)
	/*
		allAgs, agsMR, _ := Client.GetDataForFirstScreen("PB", 2018)
		agsJSON, err := json.MarshalIndent(allAgs, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n\n\n", agsJSON)
		agsJSON, err = json.MarshalIndent(agsMR, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", agsJSON)
	*/
	c.Disconnect()
}
