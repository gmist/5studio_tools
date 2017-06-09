package main

type Image struct {
	ImageURL string `json:"image_url"`
}

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
}

type ProductResponse struct {
	Status  string
	Count   int
	Now     string
	NextURL string `json:"next_url"`
	Result  []Product
}

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
