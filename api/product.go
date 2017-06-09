package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gmist/5studio_tools/lib"
)

type Product struct {
	ID              uint64
	Name            string
	Barcode         string
	Price           float64
	URL             string
	Leftovers       int         `json:"leftovers"`
	Category        string      `json:"category_name"`
	Subcategory     string      `json:"subcategory_name"`
	CountryOfOrigin string      `json:"country"`
	Vendor          string      `json:"brand"`
	Description     string      `json:"description"`
	SalesNotes      string      `json:"equipment"`
	Available       bool        `json:"is_available"`
	VendorCode      string      `json:"catalogue_id"`
	Pictures        []lib.Image `json:"images"`
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
