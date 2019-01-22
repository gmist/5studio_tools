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
		var tmp []api.Image
		for _, img := range product.Pictures {
			if img.ServeURL == "" {
				continue
			}
			tmp = append(tmp, img)
		}
		for i, img := range tmp {
			var image NameImage

			var ext string
			switch sw := img.ContentType; sw {
			case "image/jpeg":
				ext = "jpg"
			case "image/png":
				ext = "png"
			case "image/gif":
				ext = "gif"
			case "image/bpm":
				ext = "bmp"
			case "image/x-icon":
				ext = "ico"
			case "image/tiff":
				ext = ".tiff"
			case "image/svg+xml":
				ext = ".svg"
			case "image/webp":
				ext = ".webp"
			default:
				ext = "jpg"
			}

			image.Name = fmt.Sprintf("%s_%s.%s", product.ID1C, strconv.Itoa(i), ext)
			image.URL = img.ServeURL
			images = append(images, image)
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
