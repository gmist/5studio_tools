package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alehano/goyml"
)

func hash(s string) uint32 {
	if s == "" {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func getCatrories(URL string) ([]Category, string) {
	fmt.Println("Получение списка категорий", URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal("Ошибка при получении списка категорий по адресу:", URL, err.Error())
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var data categoryResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal("Ошибка декодирования списка категорий по адресу:", URL, err.Error())
	}
	var categories []Category
	for _, catJSON := range data.Result {
		categories = append(categories, Category{Name: catJSON.Name, ID: hash(catJSON.ID), Parent: hash(catJSON.Parent)})
	}
	return categories, data.NextURL
}

func getProducts(URL string) ([]Product, string) {
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

func main() {
	fmt.Println("Экспорт товаров компании \"Город Игр\" в YML-файл")
	currentTime := time.Now().Format("2006-01-02-15-04-05")

	_, err := os.Stat(YmlDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(YmlDir, 0755)
		if err != nil {
			log.Fatal("Невозможно создать директорию для экспорта", err.Error())
		}
	}
	fileName := fmt.Sprintf("%s.%s", currentTime, "yml")
	fileName = filepath.Join(YmlDir, fileName)

	var categories []Category
	url := CategoriesURL
	for {
		categoriesChunk, nextURL := getCatrories(url)
		categories = append(categories, categoriesChunk...)
		fmt.Printf("Скачено %v категорий\n", len(categories))
		if nextURL == "" {
			break
		}
		url = nextURL
		time.Sleep(500)
	}

	catMap := make(map[string]map[string]uint32, len(categories))
	for _, cat := range categories {
		catMap[cat.Name] = map[string]uint32{"id": cat.ID, "parent": cat.Parent}
	}

	var products []Product
	url = ProductsURL
	for {
		productsChunk, nextURL := getProducts(url)
		products = append(products, productsChunk...)
		fmt.Printf("Скачено %v товаров\n", len(products))
		if nextURL == "" {
			break
		}
		url = nextURL
		time.Sleep(1000)
	}

	fmt.Println("Скачивание продуктов завершено, получено", len(products), "позиций")
	fmt.Println("Генерация YML файла")

	ymlCat := goyml.NewYML("Город Игр", "Компания Город Игр", "http://5studio.ru/")
	ymlCat.Shop.Email = "i@5studio.ru"
	ymlCat.AddCurrency("RUR", "1", 0)

	for _, cat := range categories {
		if cat.Parent == 0 {
			ymlCat.AddCategory(int(cat.ID), int(cat.Parent), cat.Name)
		}
	}

	for _, cat := range categories {
		if cat.Parent != 0 {
			ymlCat.AddCategory(int(cat.ID), int(cat.Parent), cat.Name)
		}
	}

	for _, product := range products {
		if product.Leftovers <= 0 {
			continue
		}

		var categoryID uint32
		if product.Subcategory != "" {
			categoryID = catMap[product.Subcategory]["id"]
		} else {
			categoryID = catMap[product.Category]["id"]
		}
		offer := goyml.Offer{
			Id:              strconv.FormatUint(product.ID, 10),
			Name:            product.Name,
			Available:       product.Available,
			CategoryId:      int(categoryID),
			CountryOfOrigin: product.CountryOfOrigin,
			CurrencyId:      "RUR",
			Description:     product.Description,
			Price:           product.Price,
			SalesNotes:      product.SalesNotes,
			Vendor:          product.Vendor,
			VendorCode:      product.VendorCode,
			Url:             product.URL,
		}
		offer.AddBarcode(product.Barcode)
		if len(product.Pictures) > 0 {
			for _, img := range product.Pictures {
				offer.AddPicture(img.ImageURL)
			}
		}
		ymlCat.AddOffer(offer)
	}
	goyml.ExportToFile(ymlCat, fileName, true)
}
