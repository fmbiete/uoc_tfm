package orm

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const errMsgTxCommit string = "Failed to commit changes to database for order conversion"

func (d *Database) OrderLineCreate(orderId uint64, lineCart CartLine) (Order, error) {
	var err error

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(orderId)
	}

	// Convert line cart to line order
	var line OrderLine = OrderLine{Quantity: lineCart.Quantity, OrderID: orderId}

	var dish Dish
	err = d.db.Select("nombre").First(&dish, lineCart.DishID).Error
	if err != nil {
		log.Error().Err(err).Interface("line", line).Msg("Failed to read dish from cart line")
		return d.OrderDetails(orderId)
	}
	line.Name = dish.Name

	var costUnit float64
	costUnit, err = d.dishCurrentCost(lineCart.DishID)
	if err != nil {
		log.Error().Err(err).Uint64("dishId", lineCart.DishID).Msg("Failed to read cost from dish")
		return d.OrderDetails(orderId)
	}
	line.CostUnit = costUnit

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		err = tx.Save(&line).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to save line order")
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

func (d *Database) OrderLineDelete(orderId uint64, lineId uint64) (Order, error) {
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

		err = tx.Delete(&OrderLine{}, lineId).Error
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

func (d *Database) OrderLineModify(line OrderLine) (Order, error) {
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
		err = tx.Model(&line).Update("cantidad", line.Quantity).Error
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

	// Get subvention
	var config Configuration
	err := tx.Select("subvention").First(&config).Error
	if err != nil {
		log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read subvention")
		return err
	}

	// Get lines
	var lines []OrderLine
	err = tx.Select("cost_unit", "quantity").Where("order_id = ?", orderId).Find(&lines).Error
	if err != nil {
		log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read order")
		return err
	}

	// Calculate cost
	var costTotal, costToPay float64 = d.orderCalculateCostNoDB(lines, config.Subvention)

	// Update Order cost values
	var curOrder Order
	curOrder.ID = orderId
	err = tx.Model(&curOrder).Updates(Order{CostTotal: costTotal, CostToPay: costToPay}).Error
	if err != nil {
		log.Error().Err(err).Uint64("orderId", orderId).Msg("Failed to update order")
		return err
	}

	return nil
}
