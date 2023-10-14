package orm

import "github.com/rs/zerolog/log"

func (d *Database) OrderCreateFromCart(userId uint64) (Order, error) {
	var err error

	// is the kitchen open?
	err = d.configChangesAllowed()
	if err != nil {
		return Order{}, err
	}

	// get cart
	cart, err := d.CartDetails(userId)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Msg("Failed to read cart")
		return Order{}, err
	}

	// initialize order
	var order Order
	order.UserID = userId
	order.Delivery, err = d.configTodayDelivery()
	if err != nil {
		return Order{}, err
	}

	// Copy cart lines to order lines
	var dish Dish
	var costUnit float64
	for _, line := range cart.CartLines {
		err = d.db.Select("name").First(&dish, line.DishID).Error
		if err != nil {
			log.Error().Err(err).Interface("line", line).Msg("Failed to read dish from cart line")
			return Order{}, err
		}
		costUnit, err = d.dishCurrentCost(line.DishID)
		if err != nil {
			log.Error().Err(err).Uint64("dishId", line.DishID).Msg("Failed to read cost from dish")
			return Order{}, err
		}

		order.OrderLines = append(order.OrderLines, OrderLine{Name: dish.Name, Quantity: line.Quantity, CostUnit: costUnit})
	}

	// calculate order total
	order, err = d.orderCalculateCost(order)
	if err != nil {
		log.Error().Err(err).Interface("order", order).Msg("Failed to calculate cost order")
		return Order{}, err
	}

	// save order
	err = d.db.Save(&order).Error
	if err != nil {
		log.Error().Err(err).Interface("order", order).Msg("Failed to create order")
		return Order{}, err
	}

	// delete cart
	_, err = d.CartDelete(userId)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Msg("Failed to delete cart after order conversion")
		// TODO: rollback?
	}

	return order, err
}

func (d *Database) OrderDelete(orderId uint64) error {
	return d.db.Delete(&Order{}, orderId).Error
}

func (d *Database) OrderDetails(orderId uint64) (Order, error) {
	var order Order
	err := d.db.Preload("OrderLines").First(&order, orderId).Error
	return order, err
}

func (d *Database) OrderList(userId int64, day string) ([]Order, error) {
	var orders []Order

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

func (d *Database) orderCalculateCost(order Order) (Order, error) {
	// Allowed: multiple orders per day
	// TODO: apply subvention only in the first order
	subvention, err := d.configSubvention()
	if err != nil {
		return order, err
	}

	// Calculate total
	order.CostTotal = 0
	for i, _ := range order.OrderLines {
		order.CostTotal += float64(order.OrderLines[i].Quantity) * order.OrderLines[i].CostUnit
	}
	// Apply subvention
	order.CostToPay = order.CostTotal - subvention
	if order.CostToPay < 0 {
		order.CostToPay = 0
	}

	return order, nil
}

func (d *Database) orderUpdateCost(orderId uint64) (Order, error) {
	order, err := d.OrderDetails(orderId)
	if err != nil {
		log.Error().Err(err).Uint64("id", orderId).Msg("Failed to read order")
		return order, err
	}

	order, err = d.orderCalculateCost(order)
	if err != nil {
		log.Error().Err(err).Interface("order", order).Msg("Failed to calculate cost order")
		return order, err
	}

	var curOrder Order
	curOrder.ID = orderId
	err = d.db.Model(&curOrder).Updates(Order{CostTotal: order.CostTotal, CostToPay: order.CostToPay}).Error
	if err != nil {
		log.Error().Err(err).Interface("order", order).Msg("Failed to update order")
		return d.OrderDetails(orderId)
	}

	return order, nil
}
