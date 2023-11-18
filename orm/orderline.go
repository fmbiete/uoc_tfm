package orm

import (
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const errMsgTxCommit string = "Failed to commit changes to database for order conversion"

func (d *Database) OrderLineCreate(userId uint64, orderId uint64, lineOrder models.OrderLine) (models.Order, error) {
	var err error

	// Only the owner of the order can add lines to it
	canProceed, err := d.orderOwnedByUser(userId, orderId)
	if !canProceed {
		return models.Order{}, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(int64(userId), orderId)
	}

	lineOrder.OrderID = orderId

	// Overwrite Name and CostUnit (anti-tampering protection)
	var dish models.Dish
	err = d.db.Select("name").First(&dish, lineOrder.DishID).Error
	if err != nil {
		log.Error().Err(err).Interface("line", lineOrder).Msg("Failed to read dish from cart line")
		return d.OrderDetails(int64(userId), orderId)
	}
	lineOrder.Name = dish.Name

	var costUnit float64
	costUnit, err = d.dishCurrentCost(lineOrder.DishID)
	if err != nil {
		log.Error().Err(err).Uint64("dishId", lineOrder.DishID).Msg("Failed to read cost from dish")
		return d.OrderDetails(int64(userId), orderId)
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
			return d.OrderDetails(int64(userId), orderId)
		}

		err = d.orderUpdateCost(tx, orderId)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), orderId)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", orderId).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), orderId)
		}
	}

	return d.OrderDetails(int64(userId), orderId)
}

func (d *Database) OrderLineDelete(userId uint64, orderId uint64, lineId uint64) (models.Order, error) {
	var err error

	// Only the owner of the order can delete lines to it
	canProceed, err := d.orderOwnedByUser(userId, orderId)
	if !canProceed {
		return models.Order{}, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(int64(userId), orderId)
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
			return d.OrderDetails(int64(userId), orderId)
		}

		err = d.orderUpdateCost(tx, orderId)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), orderId)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", orderId).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), orderId)
		}
	}

	return d.OrderDetails(int64(userId), orderId)
}

func (d *Database) OrderLineModify(userId uint64, line models.OrderLine) (models.Order, error) {
	var err error

	// Only the owner of the order can add lines to it
	canProceed, err := d.orderOwnedByUser(userId, line.OrderID)
	if !canProceed {
		return models.Order{}, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return d.OrderDetails(int64(userId), line.OrderID)
	}

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		// existing line - update quantity ONLY
		err = tx.Model(&line).Update("quantity", line.Quantity).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to save order line")
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), line.OrderID)
		}

		// update costs in Order
		err = d.orderUpdateCost(tx, line.OrderID)
		if err != nil {
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), line.OrderID)
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("orderId", line.OrderID).Msg(errMsgTxCommit)
			// explicit rollback - we will call a non-tx function
			tx.Rollback()
			return d.OrderDetails(int64(userId), line.OrderID)
		}
	}

	return d.OrderDetails(int64(userId), line.OrderID)
}

func (d *Database) orderUpdateCost(tx *gorm.DB, orderId uint64) error {
	// Anything in this function needs to run using the transaction
	var err error
	var subvention float64 = 0

	var aux models.Order
	err = tx.Select("subvention").Where("id = ?", orderId).Find(&aux).Error
	if err != nil {
		log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read order subvention")
		return err
	}
	subvention = aux.Subvention

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
