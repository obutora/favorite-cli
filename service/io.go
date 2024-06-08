package service

import (
	"encoding/json"
	"log"
	"os"
)

type Item struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func ReadItems() []Item{
	var items []Item
	bytes, err := os.ReadFile("assets/data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		log.Fatal(err)
	}
	return items
}

func SaveItem(title string, desc string) {
	var items []Item
	bytes, err := os.ReadFile("assets/data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &items)
	if err != nil {
		log.Fatal(err)
	}

	items = append(items, Item{
		Title: title,
		Desc:  desc,
	})

	bytes, err = json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("assets/data.json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}