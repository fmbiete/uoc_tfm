package models

type CountOrders struct {
	Day   string `json:"day"`
	Count uint64 `json:"count"`
}
