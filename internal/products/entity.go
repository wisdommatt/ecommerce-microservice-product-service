package products

import "time"

type Product struct {
	ID          int       `json:"_" gorm:"autoIncrement,primaryKey"`
	Sku         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Brand       string    `json:"brand"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"imageUrl"`
	TimeAdded   time.Time `json:"timeAdded"`
}
