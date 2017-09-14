package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alehano/goyml"
	"github.com/gmist/5studio_tools/api"
	"github.com/gmist/5studio_tools/lib"
)

func main() {
	fmt.Println("Экспорт товаров компании \"Город Игр\" в YML-файл")
	currentTime := time.Now().Format("2006-01-02-15-04-05")

	_, err := os.Stat(lib.YmlDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(lib.YmlDir, 0755)
		if err != nil {
			log.Fatal("Невозможно создать директорию для экспорта", err.Error())
		}
	}
	fileName := fmt.Sprintf("%s.%s", currentTime, "xml")
	fileName = filepath.Join(lib.YmlDir, fileName)

	var categories []api.Category
	url := lib.CategoriesURL
	for {
		categoriesChunk, nextURL := api.GetCatrories(url)
		categories = append(categories, categoriesChunk...)
		fmt.Printf("Скачено %v категорий\n", len(categories))
		if nextURL == "" {
			break
		}
		url = nextURL
	}

	catMap := make(map[string]map[string]uint32, len(categories))
	for _, cat := range categories {
		catMap[strings.ToLower(strings.TrimSpace(cat.Name))] = map[string]uint32{"id": cat.ID, "parent": cat.Parent}
	}

	var products []api.Product
	url = lib.ProductsURL
	for {
		productsChunk, nextURL := api.GetProducts(url)
		products = append(products, productsChunk...)
		fmt.Printf("Скачено %v товаров\n", len(products))
		if nextURL == "" {
			break
		}
		url = nextURL
	}

	fmt.Println("Скачивание продуктов завершено, получено", len(products), "позиций")
	fmt.Println("Генерация YML файла")

	ymlCat := goyml.NewYML("Город Игр", "Компания Город Игр", "http://5studio.ru/")
	ymlCat.Shop.Email = "i@5studio.ru"
	ymlCat.AddCurrency("RUR", "1", 0)

	for _, cat := range categories {
		if cat.Parent == 0 {
			ymlCat.AddCategory(int(cat.ID), int(cat.Parent), strings.TrimSpace(cat.Name))
		}
	}

	for _, cat := range categories {
		if cat.Parent != 0 {
			ymlCat.AddCategory(int(cat.ID), int(cat.Parent), strings.TrimSpace(cat.Name))
		}
	}

	for _, product := range products {
		if product.Leftovers <= 0 {
			continue
		}

		var categoryID uint32
		if product.Subcategory != "" {
			categoryID = catMap[strings.ToLower(strings.TrimSpace(product.Subcategory))]["id"]
		} else if product.Category != "" {
			categoryID = catMap[strings.ToLower(strings.TrimSpace(product.Category))]["id"]
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
			sort.Sort(api.ByPriority(product.Pictures))
			for _, img := range product.Pictures {
				offer.AddPicture(img.URL)
			}
		}
		ymlCat.AddOffer(offer)
	}
	goyml.ExportToFile(ymlCat, fileName, true)
}
