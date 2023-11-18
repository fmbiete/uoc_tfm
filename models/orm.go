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
}

type Category struct {
	BaseModel
	Name string `gorm:"uniqueIndex;size:250"`
}

type Ingredient struct {
	BaseModel
	Name string `gorm:"uniqueIndex;size:250"`
}

type Allergen struct {
	BaseModel
	Name string `gorm:"uniqueIndex;size:250"`
}

type Promotion struct {
	BaseModel
	DishID    uint64 // FK - Promotion belongs to Dish
	Dish      Dish   // for preload joins (not reflected in model)
	StartTime time.Time
	EndTime   time.Time
	Cost      float64 `gorm:"scale:2"`
}

type Dish struct {
	BaseModel
	Name        string       `gorm:"unique;size:250"`
	Description string       `gorm:"size:2000"`
	Categories  []Category   `gorm:"many2many:dish_categories;"`
	Ingredients []Ingredient `gorm:"many2many:dish_ingredients;"`
	Allergens   []Allergen   `gorm:"many2many:dish_allergens;"`
	Cost        float64      `gorm:"scale:2"`
	Promotions  []Promotion  // has many
	Likes       uint64       `gorm:"default:0"`
	Dislikes    uint64       `gorm:"default:0"`
}

type DishLike struct {
	BaseModel
	DishID uint64 `gorm:"uniqueIndex:ix_user_like;"` // FK
	UserID uint64 `gorm:"uniqueIndex:ix_user_like;"` // FK
}

type DishDislike struct {
	BaseModel
	DishID uint64 `gorm:"uniqueIndex:ix_user_dislike;"` // FK
	UserID uint64 `gorm:"uniqueIndex:ix_user_dislike;"` // FK
}

type OrderLine struct {
	BaseModel
	OrderID  uint64  // FK - line belongs to Order
	DishID   uint64  // FK - line has 1 Dish
	Name     string  `gorm:"size:250"` // don't use dish references - attributes will change
	CostUnit float64 `gorm:"scale:2"`  // don't use dish references - attributes will change
	Quantity uint
}

type Order struct {
	BaseModel
	OrderLines    []OrderLine
	UserID        uint64  // FK - Order belongs to User
	User          User    // For preload joins, not reflected in model
	CostTotal     float64 `gorm:"scale:2"`
	CostToPay     float64 `gorm:"scale:2"` // cost to pay after subvention
	Subvention    float64 `gorm:"scale:2"` // subvention applied
	Delivery      time.Time
	Address1      string `gorm:"size:250"`
	Address2      string `gorm:"size:250"`
	Address3      string `gorm:"size:250"`
	City          string `gorm:"size:250"`
	PostalCode    string `gorm:"size:10"`
	Phone         string `gorm:"size:20"`
	PaymentMethod string `gorm:"size:250"`
	PaymentSecret string `gorm:"size:250"`
}
