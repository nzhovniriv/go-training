package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	Firstname  string `json:"first"`            // should be encoded as 'first'
	Middlename string `json:"middle,omitempty"` // should be encoded as 'middle', and not present if blank
	Lastname   string `json:"last"`             // should be encoded as 'last'

	SSID int64 `json:"-"` // should not be encoded

	City    string `json:"city,omitempty"` // should be encoded as 'city' and not present if missing
	Country string `json:"country"`        // should be encoded as 'country'

	Telephone int64 `json:"tel,string"` // should be encoded as 'tel', the value should be a string, not a number
}

var persons = []Person{
	{Firstname: "Peter", Middlename: "Jone", Lastname: "White", SSID: 1794, City: "Tokyo", Country: "Japan", Telephone: 813},
	{Firstname: "Bob", Middlename: "Mel", Lastname: "Lee", SSID: 3254, Country: "USA", Telephone: 617},
	{Firstname: "Angela", Middlename: "", Lastname: "Hill", SSID: 3657, City: "New York", Country: "USA", Telephone: 646},
	{Firstname: "Marty", Middlename: "Ado", Lastname: "Wood", SSID: 6524, City: "Las Vegas", Country: "USA", Telephone: 702},
}

func main() {
	data, err := json.Marshal(persons)
	// data, err := json.MarshalIndent(persons, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}
