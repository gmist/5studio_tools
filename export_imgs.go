package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/gmist/5studio_tools/api"
	"github.com/gmist/5studio_tools/lib"
)

type NameImage struct {
	Name string
	URL  string
}

func downloadImage(image NameImage) {
	fmt.Printf("Скачиваем изображение %s (%s)\n", image.Name, image.URL)
	resp, err := http.Get(image.URL)
	if err != nil {
		log.Fatalf("Ошибка получения изображения по адресу: %s -- %s", image.URL, err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения изображения по адресу: %s -- %s", image.URL, err)
	}

	filename := path.Join(lib.ImgsDir, image.Name)
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Ошибка создания фалов: %s -- %s", image.URL, err)
		}
		file.Close()
		err = ioutil.WriteFile(filename, contents, 0644)
		if err != nil {
			log.Fatalf("Ошибка записи в файл: %s -- %s", filename, err)
		}
	}
}

func main() {
	fmt.Println("Экспорт основных изображений")

	_, err := os.Stat(lib.ImgsDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(lib.ImgsDir, 0755)
		if err != nil {
			log.Fatal("Невозможно создать директорию для экспорта изображений", err.Error())
		}
	}
	var products []api.Product
	url := lib.ProductsURL
	for {
		if url == "" {
			break
		}
		productsChunk, nextURL := api.GetProducts(url)
		products = append(products, productsChunk...)
		fmt.Printf("Скачено %v товаров\n", len(products))
		url = nextURL
	}

	fmt.Println("Скачивание продуктов завершено, получено", len(products), "позиций")
	fmt.Println("Приступаем к скачиванию изображений")

	var images []NameImage
	for _, product := range products {
		var tmp [3]api.Image
		for _, img := range product.Pictures {
			if img.URL == "" {
				continue
			}
			if img.Type == 70 && tmp[0].Type == 0 {
				tmp[0] = img
			} else if img.Type == 10 && tmp[1].Type == 0 {
				tmp[1] = img
			} else if img.Type == 20 && tmp[2].Type == 0 {
				tmp[2] = img
			}
		}
		for i, img := range tmp {
			if img.Type != 0 {
				var image NameImage
				image.Name = fmt.Sprintf("%s_%s.jpg", product.ID1C, strconv.Itoa(i))
				image.URL = fmt.Sprintf("%s=s100", img.URL)
				images = append(images, image)
			}
		}
	}
	var throttle = make(chan int, 5)
	var wg sync.WaitGroup
	fmt.Println("Найдено", len(images), "изображений")
	wg.Add(len(images))
	for _, image := range images {
		throttle <- 1
		go func(image NameImage) {
			defer wg.Done()
			downloadImage(image)
			<-throttle
		}(image)
	}
	wg.Wait()
}
