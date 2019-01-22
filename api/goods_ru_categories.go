package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GoodsRuCategory struct {
	Name      string
	FullName  string `json:"full_name"`
	ID        uint64
	ParentKey string `json:"parent_key,omitempty"`
	Key       string
}

type GoodsRuCategoryResponse struct {
	Status  string
	Count   int
	Now     string
	NextURL string `json:"next_url"`
	Result  []GoodsRuCategory
}

// type GoodsRuCategory struct {
// 	Name      string
// 	ParentKey string
// 	ID        uint64
// 	FullName  string
// 	Key       string
// }

func GetGoodsRuCatrories(URL string) ([]GoodsRuCategory, string) {
	fmt.Println("Получение списка категорий", URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Ошибка при получении списка категорий по адресу: ", URL, " ", err.Error())
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data GoodsRuCategoryResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal("Ошибка декодирования списка категорий по адресу: ", URL, " ", err.Error())
	}
	// var categories []GoodsRuCategory
	// for _, catJSON := range data.Result {
	// 	categories = append(
	// 		categories,
	// 		GoodsRuCategory{Name: catJSON.Name, Key: catJSON.Key, ParentKey: catJSON.ParentKey, FullName: catJSON.FullName, ID: catJSON.ID})
	// }
	return data.Result, data.NextURL
}
