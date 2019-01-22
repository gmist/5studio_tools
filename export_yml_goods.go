package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gmist/5studio_tools/api"
	"github.com/gmist/5studio_tools/lib"
	"github.com/gmist/goyml"
)

func main() {
	fmt.Println("Экспорт товаров компании \"Город Игр\" в YML-файл для goods.ru")
	currentTime := time.Now().Format("2006-01-02-15-04-05")

	var categories []api.GoodsRuCategory
	url := lib.GoodsRuCategoriesURL
	for {
		categoriesChunk, nextURL := api.GetGoodsRuCatrories(url)
		categories = append(categories, categoriesChunk...)
		fmt.Printf("Скачено %v категорий\n", len(categories))
		if nextURL == "" {
			break
		}
		url = nextURL
	}

	if len(categories) == 0 {
		fmt.Println("Категории не найдены - ошибка: недопустимо отсуствие категорий")
		fmt.Println("Press any key to exit")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	fmt.Println("\nСписок категорий goods.ru:")
	for _, category := range categories {
		fmt.Println(category.FullName)
	}

	var products []api.Product
	url = lib.GoodsRuProductsURL
	for {
		productsChunk, nextURL := api.GetProducts(url)
		products = append(products, productsChunk...)
		fmt.Printf("Скачено %v товара(ов)\n", len(products))
		if nextURL == "" {
			break
		}
		url = nextURL
	}

	if len(products) == 0 {
		fmt.Println("Продукты не найдены - ошибка: недопустимо осутствие товаров в выгрузке")
		fmt.Println("Press any key to exit")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	fmt.Println("Скачивание продуктов завершено, получено", len(products), "позиций")
	fmt.Println("\nНачинаем генерацию XML файла")

	categoriesMap := make(map[string]api.GoodsRuCategory)
	for _, category := range categories {
		categoriesMap[category.Key] = category
	}

	mergeCategoriesMap := make(map[string]api.GoodsRuCategory)
	for _, product := range products {
		tmpKey := product.GoodsCategory.Key
		mergeCategoriesMap[tmpKey] = product.GoodsCategory
		for {
			if tmpKey == "" {
				break
			}
			mergeCategoriesMap[tmpKey] = categoriesMap[tmpKey]
			tmpKey = mergeCategoriesMap[tmpKey].ParentKey
		}
	}

	var sortedCategoriesNames []string
	sortedCategories := map[string]api.GoodsRuCategory{}
	for k, v := range mergeCategoriesMap {
		sortedCategories[v.FullName] = mergeCategoriesMap[k]
		sortedCategoriesNames = append(sortedCategoriesNames, v.FullName)
	}
	sort.Sort(sort.StringSlice(sortedCategoriesNames))

	ymlCat := goyml.NewYML("Город Игр", `Компания "Город Игр"`, "http://www.5studio.ru/")
	ymlCat.AddCurrency("RUR", "1", 0)

	fmt.Println("\nАнализируем категории товаров")
	for _, v := range sortedCategoriesNames {
		tmpCat := sortedCategories[v]
		parentID := uint64(0)
		if tmpCat.ParentKey != "" {
			parentID = mergeCategoriesMap[tmpCat.ParentKey].ID
		}
		ymlCat.AddCategory(tmpCat.ID, parentID, tmpCat.Name)
	}

	fmt.Println("Добавляем товары в выгрузку")
	for _, product := range products {
		if product.Leftovers <= 0 || product.Available == false {
			fmt.Println(product.Name, "- пропущен по причине отсутствия")
			continue
		}

		offer := goyml.Offer{
			Id:              strconv.FormatUint(product.ID, 10),
			Name:            product.Name,
			Available:       product.Available,
			CategoryId:      product.GoodsCategory.ID,
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
		fmt.Println(product.Name, "- добавлен")
	}

	fmt.Println("\nЗаписываем файл на диск")
	fileName := fmt.Sprintf("%s.xml", currentTime)
	goyml.ExportToFile(ymlCat, fileName, true)

	fmt.Println("Генерация XML файла закончена:", fileName)
	fmt.Println("\nPress any key to exit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
