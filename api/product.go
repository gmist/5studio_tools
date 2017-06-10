package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Image struct {
	URL      string `json:"image_url"`
	Type     int    `json:"image_type"`
	Priority int    `json:"priority"`
}

type ByPriority []Image

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority > a[j].Priority }

type Product struct {
	ID              uint64
	Name            string
	Barcode         string
	Price           float64
	URL             string
	Leftovers       int     `json:"leftovers"`
	Category        string  `json:"category_name"`
	Subcategory     string  `json:"subcategory_name"`
	CountryOfOrigin string  `json:"country"`
	Vendor          string  `json:"brand"`
	Description     string  `json:"description"`
	SalesNotes      string  `json:"equipment"`
	Available       bool    `json:"is_available"`
	VendorCode      string  `json:"catalogue_id"`
	Pictures        []Image `json:"images"`
	ID1C            string  `json:"id_1c"`
}

type ProductResponse struct {
	Status  string
	Count   int
	Now     string
	NextURL string `json:"next_url"`
	Result  []Product
}

func GetProducts(URL string) ([]Product, string) {
	fmt.Println("Получение списка товаров", URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Ошибка при получении списка товаров по адресу:", URL, err.Error())
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data ProductResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal("Ошибка декодирования списка товаров по адресу:", URL, err.Error())
	}
	return data.Result, data.NextURL
}
