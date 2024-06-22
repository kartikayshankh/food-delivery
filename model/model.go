package model

type User struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
	Password    string `json:"password" bson:"password"`
}

type Rider struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Email       string   `json:"email" bson:"email"`
	Password    string   `json:"password" bson:"password"`
	PhoneNumber string   `json:"phone_number" bson:"phone_number"`
	VehicleType string   `json:"vehicle_type" bson:"vehicle_type"`
	Location    Location `json:"location" bson:"location"`
	Orders      []Order  `json:"orders" bson:"orders"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Restaurant struct {
	ID      string     `json:"id" bson:"_id,omitempty"`
	Name    string     `json:"name" bson:"name"`
	Address string     `json:"address" bson:"address"`
	Cuisine string     `json:"cuisine" bson:"cuisine"`
	Menu    []MenuItem `json:"menu" bson:"menu"`
	Rating  float64    `json:"rating" bson:"rating"`
	Orders  []Order    `json:"orders" bson:"orders"`
}

type MenuItem struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Price       float64 `json:"price" bson:"price"`
}

type Order struct {
	UserID     string      `json:"user_id" bson:"user_id"`
	Items      []OrderItem `json:"items" bson:"items"`
	TotalPrice float64     `json:"total_price" bson:"total_price"`
	Status     string      `json:"status" bson:"status"`
}

type OrderItem struct {
	MenuItemID string `json:"menu_item_id" bson:"menu_item_id"`
	Quantity   int    `json:"quantity" bson:"quantity"`
}
