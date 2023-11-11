package orm

import (
	"errors"
	"tfm_backend/models"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) OrderCreate(order models.Order) (models.Order, error) {
	var err error

	// is the kitchen open?
	err = d.configChangesAllowed()
	if err != nil {
		return models.Order{}, err
	}

	// initialize order
	order.Delivery, err = d.configTodayDelivery()
	if err != nil {
		return models.Order{}, err
	}

	// Overwrite order line prices with dish prices (tampering protection)
	var dish models.Dish
	var costUnit float64
	for i, line := range order.OrderLines {
		err = d.db.Select("name").First(&dish, line.DishID).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to read dish from order line")
			return models.Order{}, err
		}
		costUnit, err = d.dishCurrentCost(line.DishID)
		if err != nil {
			log.Error().Err(err).Uint64("dishId", line.DishID).Msg("Failed to read cost from dish")
			return models.Order{}, err
		}

		order.OrderLines[i].CostUnit = costUnit
		order.OrderLines[i].Name = dish.Name
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

		// create order and lines
		err = tx.Create(&order).Error
		if err != nil {
			log.Error().Err(err).Interface("order", order).Msg("Failed to create order")
			return models.Order{}, err
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Interface("order", order).Msg(errMsgTxCommit)
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
	err := d.db.Preload("OrderLines", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_lines.name")
	}).First(&order, orderId).Error
	return order, err
}

func (d *Database) OrderList(userId int64, day string, limit uint64, offset uint64) ([]models.Order, error) {
	var orders []models.Order

	queryDb := d.db.Where("1 = 1")

	if len(day) > 0 {
		queryDb = queryDb.Where("date(created_at) = ?", day)
	}
	if userId > 0 {
		queryDb = queryDb.Where("user_id = ?", userId)
	}

	err := queryDb.Limit(int(limit)).Offset(int(offset)).Find(&orders).Error
	return orders, err
}

func (d *Database) OrderSubvention(userId uint64) (float64, error) {
	var err error
	var subvention float64 = 0

	// is the kitchen open?
	err = d.configChangesAllowed()
	if err != nil {
		return subvention, err
	}

	subvention, err = d.orderCalculateSubvention(userId)

	return subvention, err
}

func (d *Database) orderCalculateCost(order models.Order) (models.Order, error) {
	var err error
	var subvention float64 = 0

	subvention, err = d.orderCalculateSubvention(order.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate cost - calculate subvention")
		return order, err
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

func (d *Database) orderCalculateSubvention(userId uint64) (float64, error) {
	var err error
	var subvention float64 = 0

	// Allowed: multiple orders per day, but only first has subvention
	err = d.db.Where("date(created_at) = date(?) AND user_id = ?", time.Now(), userId).First(&models.Order{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Get subvention - this is the first order
		subvention, err = d.configSubvention()
	}

	return subvention, err
}
