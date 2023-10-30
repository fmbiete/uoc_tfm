package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const errMsgTxCommit string = "Failed to commit changes to database for order conversion"

func (d *Database) OrderLineCreate(orderId uint64, lineOrder models.OrderLine) (models.Order, error) {
	var err error

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(orderId)
	}

	lineOrder.OrderID = orderId

	// Overwrite Name and CostUnit (anti-tampering protection)
	var dish models.Dish
	err = d.db.Select("name").First(&dish, lineOrder.DishID).Error
	if err != nil {
		log.Error().Err(err).Interface("line", lineOrder).Msg("Failed to read dish from cart line")
		return d.OrderDetails(orderId)
	}
	lineOrder.Name = dish.Name

	var costUnit float64
	costUnit, err = d.dishCurrentCost(lineOrder.DishID)
	if err != nil {
		log.Error().Err(err).Uint64("dishId", lineOrder.DishID).Msg("Failed to read cost from dish")
		return d.OrderDetails(orderId)
	}
	lineOrder.CostUnit = costUnit

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		err = tx.Save(&lineOrder).Error
		if err != nil {
			log.Error().Err(err).Interface("line", lineOrder).Msg("Failed to save line order")
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}

		err = d.orderUpdateCost(tx, orderId)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", orderId).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}
	}

	return d.OrderDetails(orderId)
}

func (d *Database) OrderLineDelete(orderId uint64, lineId uint64) (models.Order, error) {
	var err error

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(orderId)
	}

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		err = tx.Delete(&models.OrderLine{}, lineId).Error
		if err != nil {
			log.Error().Err(err).Uint64("lineId", lineId).Msg("Failed to delete line order")
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}

		err = d.orderUpdateCost(tx, orderId)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", orderId).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(orderId)
		}
	}

	return d.OrderDetails(orderId)
}

func (d *Database) OrderLineModify(line models.OrderLine) (models.Order, error) {
	var err error

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(line.OrderID)
	}

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		// existing line - update cantidad ONLY
		err = tx.Model(&line).Update("quantity", line.Quantity).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to save order line")
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(line.OrderID)
		}

		// use save to update or create line as required
		err = d.orderUpdateCost(tx, line.OrderID)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(line.OrderID)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", line.OrderID).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(line.OrderID)
		}
	}

	return d.OrderDetails(line.OrderID)
}

func (d *Database) orderUpdateCost(tx *gorm.DB, orderId uint64) error {
	// Anything in this function needs to run using the transaction
	var err error
	var subvention float64 = 0

	// Only the first order has subvention
	err = tx.Select("id").Where("id = ? AND cost_total != cost_to_pay").Find(&models.Order{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// If the costTotal != costToPay, this order has subvention applied. This is the first, so get subvention
		var config models.Configuration
		err = tx.Select("subvention").First(&config).Error
		if err != nil {
			log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read subvention")
			return err
		}
		subvention = config.Subvention
	}

	// Get lines
	var lines []models.OrderLine
	err = tx.Select("cost_unit", "quantity").Where("order_id = ?", orderId).Find(&lines).Error
	if err != nil {
		log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read order")
		return err
	}

	// Calculate cost
	var costTotal, costToPay float64 = d.orderCalculateCostNoDB(lines, subvention)

	// Update Order cost values
	var curOrder models.Order
	curOrder.ID = orderId
	err = tx.Model(&curOrder).Updates(models.Order{CostTotal: costTotal, CostToPay: costToPay}).Error
	if err != nil {
		log.Error().Err(err).Uint64("orderId", orderId).Msg("Failed to update order")
		return err
	}

	return nil
}
