package orm

import (
	"errors"
	"tfm_backend/models"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) OrderCreateFromCart(userId uint64) (models.Order, error) {
	var err error

	// is the kitchen open?
	err = d.configChangesAllowed()
	if err != nil {
		return models.Order{}, err
	}

	// get cart
	cart, err := d.CartDetails(userId)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Msg("Failed to read cart")
		return models.Order{}, err
	}

	// initialize order
	var order models.Order
	order.UserID = userId
	order.Delivery, err = d.configTodayDelivery()
	if err != nil {
		return models.Order{}, err
	}

	// Copy cart lines to order lines
	var dish models.Dish
	var costUnit float64
	for _, line := range cart.CartLines {
		err = d.db.Select("name").First(&dish, line.DishID).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to read dish from cart line")
			return models.Order{}, err
		}
		costUnit, err = d.dishCurrentCost(line.DishID)
		if err != nil {
			log.Error().Err(err).Uint64("dishId", line.DishID).Msg("Failed to read cost from dish")
			return models.Order{}, err
		}

		order.OrderLines = append(order.OrderLines, models.OrderLine{DishID: line.DishID, Name: dish.Name,
			Quantity: line.Quantity, CostUnit: costUnit})
	}

	// calculate order total
	order, err = d.orderCalculateCost(order)
	if err != nil {
		log.Error().Err(err).Interface("order", order).Msg("Failed to calculate cost order")
		return models.Order{}, err
	}

	// Transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		// save order
		err = tx.Save(&order).Error
		if err != nil {
			log.Error().Err(err).Interface("order", order).Msg("Failed to create order")
			return models.Order{}, err
		}

		// delete cart lines, we don't need to modify Cart
		err = tx.Unscoped().Model(&cart).Association("CartLines").Unscoped().Clear()
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Msg("Failed to delete cart after order conversion")
			return models.Order{}, err
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Msg(errMsgTxCommit)
			return models.Order{}, err
		}
	}

	return order, err
}

func (d *Database) OrderDelete(orderId uint64) error {
	return d.db.Delete(&models.Order{}, orderId).Error
}

func (d *Database) OrderDetails(orderId uint64) (models.Order, error) {
	var order models.Order
	err := d.db.Preload("OrderLines").First(&order, orderId).Error
	return order, err
}

func (d *Database) OrderList(userId int64, day string) ([]models.Order, error) {
	var orders []models.Order

	queryDb := d.db.Where("1 = 1")

	if len(day) > 0 {
		queryDb = queryDb.Where("date(created_at) = ?", day)
	}
	if userId > 0 {
		queryDb = queryDb.Where("user_id = ?", userId)
	}

	err := queryDb.Find(&orders).Error
	return orders, err
}

func (d *Database) orderCalculateCost(order models.Order) (models.Order, error) {
	var err error
	var subvention float64 = 0

	// Allowed: multiple orders per day, but only first has subvention
	err = d.db.Where("date(created_at) = date(?) AND id != ?", time.Now, order.ID).First(&models.Order{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Get subvention - this is the first order
		subvention, err = d.configSubvention()
		if err != nil {
			return order, err
		}
	}

	order.CostTotal, order.CostToPay = d.orderCalculateCostNoDB(order.OrderLines, subvention)

	return order, nil
}

func (d *Database) orderCalculateCostNoDB(lines []models.OrderLine, subvention float64) (float64, float64) {
	var costTotal, costToPay float64 = 0, 0

	// Calculate total
	for i := range lines {
		costTotal += float64(lines[i].Quantity) * lines[i].CostUnit
	}
	// Apply subvention
	costToPay = costTotal - subvention
	if costToPay < 0 {
		costToPay = 0
	}

	return costTotal, costToPay
}
