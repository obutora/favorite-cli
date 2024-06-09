package service

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
)

type Item struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func ReadItems() []Item {
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

func sortItemsAsc(items []Item) []Item {
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})
	return items

}

func saveItems(items []Item) {
	bytes, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("assets/data.json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func AddItem(title string, desc string) {
	items := ReadItems()

	items = append(items, Item{
		Title: title,
		Desc:  desc,
	})

	items = sortItemsAsc(items)
	saveItems(items)
}

func DeleteItem(title string) {
	items := ReadItems()

	var newItems []Item
	for _, e := range items {
		if e.Title != title {
			newItems = append(newItems, e)
		}
	}

	saveItems(newItems)
}
