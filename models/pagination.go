package models

type PaginationDishes struct {
	Dishes []Dish `json:"dishes"`
	Page   uint64 `json:"page"`
	Limit  uint64 `json:"limit"`
}

type PaginationOrders struct {
	Orders []Order `json:"orders"`
	Page   uint64  `json:"page"`
	Limit  uint64  `json:"limit"`
}

type PaginationPromotion struct {
	Promotions []Promotion `json:"promotions"`
	Page       uint64      `json:"page"`
	Limit      uint64      `json:"limit"`
}

type PaginationUsers struct {
	Users []User `json:"users"`
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
}
