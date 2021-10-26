package products

type Product struct {
	ID          int     `json:"_" gorm:"autoIncrement,primaryKey"`
	Sku         string  `json:"sku,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Category    string  `json:"category,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	Price       float64 `json:"price,omitempty"`
}
