package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gmist/5studio_tools/lib"
)

type CategoryJSON struct {
	Name   string
	Parent string `json:"root_key,omitempty"`
	ID     string `json:"key"`
}

type CategoryResponse struct {
	Status  string
	Count   int
	Now     string
	NextURL string `json:"next_url"`
	Result  []CategoryJSON
}

type Category struct {
	Name   string
	Parent uint32
	ID     uint32
}

func GetCatrories(URL string) ([]Category, string) {
	fmt.Println("Получение списка категорий", URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Ошибка при получении списка категорий по адресу:", URL, err.Error())
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data CategoryResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal("Ошибка декодирования списка категорий по адресу:", URL, err.Error())
	}
	var categories []Category
	for _, catJSON := range data.Result {
		categories = append(categories, Category{Name: catJSON.Name, ID: lib.Hash(catJSON.ID), Parent: lib.Hash(catJSON.Parent)})
	}
	return categories, data.NextURL
}
