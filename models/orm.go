package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Configuration struct {
	BaseModel
	DeliveryTime time.Time
	ChangesTime  time.Time
	Subvention   float64
}

type User struct {
	BaseModel
	Email      string `gorm:"size:100;index:ix_login,priority:1"`
	Password   string `gorm:"size:250;index:ix_login,priority:2"`
	Name       string `gorm:"size:250"`
	Surname    string `gorm:"size:250"`
	Address1   string `gorm:"size:250"`
	Address2   string `gorm:"size:250"`
	Address3   string `gorm:"size:250"`
	City       string `gorm:"size:250"`
	PostalCode string `gorm:"size:10"`
	Phone      string `gorm:"size:20"`
	IsAdmin    bool
	Orders     []Order // has many
	CartOrder  Cart    // has one
}

type Ingredient struct {
	BaseModel
	Name string `gorm:"size:250"`
}

type Allergen struct {
	BaseModel
	Name string `gorm:"size:250"`
}

type Promotion struct {
	BaseModel
	DishID    uint // FK - Promotion belongs to Dish
	StartTime time.Time
	EndTime   time.Time
	Cost      float64 `gorm:"scale:2"`
}

type Dish struct {
	BaseModel
	Name        string       `gorm:"unique;size:250"`
	Description string       `gorm:"size:2000"`
	Ingredients []Ingredient `gorm:"many2many:dish_ingredients;"`
	Allergens   []Allergen   `gorm:"many2many:dish_allergens;"`
	Cost        float64      `gorm:"scale:2"`
	Promotions  []Promotion  // has many
}

type OrderLine struct {
	BaseModel
	OrderID  uint64  // FK - line belongs to Order
	Name     string  `gorm:"size:250"` // don't use dish references - attributes will change
	CostUnit float64 `gorm:"scale:2"`
	Quantity uint
}

type Order struct {
	BaseModel
	OrderLines []OrderLine
	UserID     uint64  // FK - Order belongs to User
	CostTotal  float64 `gorm:"scale:2"`
	CostToPay  float64 `gorm:"scale:2"` // FK - cost to pay after subvention
	Delivery   time.Time
}

type CartLine struct {
	BaseModel
	CartID   uint64 // FK - CartLine belongs to Cart
	DishID   uint64 // FK - CartLine has 1 Dish
	Quantity uint
}

type Cart struct {
	BaseModel
	CartLines []CartLine
	UserID    uint64 // FK - Cart belongs to User
}
