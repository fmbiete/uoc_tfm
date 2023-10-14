package orm

import "github.com/rs/zerolog/log"

func (d *Database) OrderLineCreate(orderId uint64, lineCart CartLine) (Order, error) {
	order, err := d.OrderDetails(orderId)
	if err != nil {
		log.Error().Err(err).Uint64("orderId", orderId).Msg("Failed to read order")
		return order, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return order, err
	}

	// Convert line cart to line order
	var line OrderLine = OrderLine{Quantity: lineCart.Quantity, OrderID: orderId}

	var dish Dish
	err = d.db.Select("nombre").First(&dish, lineCart.DishID).Error
	if err != nil {
		log.Error().Err(err).Interface("line", line).Msg("Failed to read dish from cart line")
		return order, err
	}
	line.Name = dish.Name

	var costUnit float64
	costUnit, err = d.dishCurrentCost(lineCart.DishID)
	if err != nil {
		log.Error().Err(err).Uint64("dishId", lineCart.DishID).Msg("Failed to read cost from dish")
		return order, err
	}
	line.CostUnit = costUnit

	err = d.db.Save(&line).Error
	if err != nil {
		log.Error().Err(err).Interface("line", line).Msg("Failed to save line order")
		return order, err
	}

	return d.orderUpdateCost(orderId)
}

func (d *Database) OrderLineDelete(orderId uint64, lineId uint64) (Order, error) {
	var err error
	var order Order

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return order, err
	}

	err = d.db.Delete(&OrderLine{}, lineId).Error
	if err != nil {
		log.Error().Err(err).Uint64("lineId", lineId).Msg("Failed to delete line order")
		return order, err
	}

	return d.orderUpdateCost(orderId)
}
func (d *Database) OrderLineModify(line OrderLine) (Order, error) {
	order, err := d.OrderDetails(line.OrderID)
	if err != nil {
		log.Error().Err(err).Uint64("orderId", line.OrderID).Msg("Failed to read order")
		return order, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return order, err
	}

	// existing line - update cantidad ONLY
	err = d.db.Model(&line).Update("cantidad", line.Quantity).Error
	if err != nil {
		log.Error().Err(err).Interface("line", line).Msg("Failed to save order line")
		return order, err
	}

	// use save to update or create line as required
	return d.orderUpdateCost(line.OrderID)
}
