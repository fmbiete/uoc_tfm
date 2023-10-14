package orm

import "github.com/rs/zerolog/log"

func (d *Database) PedidoCreateFromCarrito(userId uint64) (Pedido, error) {
	var err error

	// is the kitchen open?
	err = d.configChangesAllowed()
	if err != nil {
		return Pedido{}, err
	}

	// get carrito
	carrito, err := d.CarritoDetails(userId)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Msg("Failed to read carrito")
		return Pedido{}, err
	}

	// initialize pedido
	var pedido Pedido
	pedido.UsuarioID = userId
	pedido.Entrega, err = d.configTodayEntrega()
	if err != nil {
		return Pedido{}, err
	}

	// Copy carrito lineas to pedido lineas
	var plato Plato
	var precioUnidad float64
	for _, linea := range carrito.CarritoLineas {
		err = d.db.Select("nombre").First(&plato, linea.PlatoID).Error
		if err != nil {
			log.Error().Err(err).Interface("linea", linea).Msg("Failed to read plato from carrito linea")
			return Pedido{}, err
		}
		precioUnidad, err = d.platoCurrentPrecio(linea.PlatoID)
		if err != nil {
			log.Error().Err(err).Uint64("platoId", linea.PlatoID).Msg("Failed to read precio from plato")
			return Pedido{}, err
		}

		pedido.PedidoLineas = append(pedido.PedidoLineas, PedidoLinea{Nombre: plato.Nombre, Cantidad: linea.Cantidad, PrecioUnidad: precioUnidad})
	}

	// calculate pedido total
	pedido, err = d.pedidoCalcularPrecio(pedido)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to calculate pedido precio")
		return Pedido{}, err
	}

	// save pedido
	err = d.db.Save(&pedido).Error
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to create pedido")
		return Pedido{}, err
	}

	// delete carrito
	_, err = d.CarritoDelete(userId)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Msg("Failed to delete carrito after pedido conversion")
		// TODO: rollback?
	}

	return pedido, err
}

func (d *Database) PedidoDelete(pedidoId uint64) error {
	return d.db.Delete(&Pedido{}, pedidoId).Error
}

func (d *Database) PedidoDetails(pedidoId uint64) (Pedido, error) {
	var pedido Pedido
	result := d.db.Preload("PedidoLineas").First(&pedido, pedidoId)
	return pedido, result.Error
}

func (d *Database) PedidoList(userId int64, day string) ([]Pedido, error) {
	var pedidos []Pedido

	queryDb := d.db.Where("1 = 1")

	if len(day) > 0 {
		queryDb = queryDb.Where("date(created_at) = ?", day)
	}
	if userId > 0 {
		queryDb = queryDb.Where("usuario_id = ?", userId)
	}

	result := queryDb.Find(&pedidos)
	return pedidos, result.Error
}

func (d *Database) pedidoCalcularPrecio(pedido Pedido) (Pedido, error) {
	// Allowed: multiple pedidos per day
	// TODO: apply subvencion only in the first pedido
	subvencion, err := d.configPrecioSubvencion()
	if err != nil {
		return pedido, err
	}

	// Calculate total
	pedido.PrecioTotal = 0
	for i, _ := range pedido.PedidoLineas {
		pedido.PrecioTotal += float64(pedido.PedidoLineas[i].Cantidad) * pedido.PedidoLineas[i].PrecioUnidad
	}
	// Apply subvencion
	pedido.PrecioPagar = pedido.PrecioTotal - subvencion
	if pedido.PrecioPagar < 0 {
		pedido.PrecioPagar = 0
	}

	return pedido, nil
}

func (d *Database) pedidoUpdatePrecio(pedidoId uint64) (Pedido, error) {
	pedido, err := d.PedidoDetails(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("id", pedidoId).Msg("Failed to read pedido")
		return pedido, err
	}

	pedido, err = d.pedidoCalcularPrecio(pedido)
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to calculate precio pedido")
		return pedido, err
	}

	var curPedido Pedido
	curPedido.ID = pedidoId
	err = d.db.Model(&curPedido).Updates(Pedido{PrecioTotal: pedido.PrecioTotal, PrecioPagar: pedido.PrecioPagar}).Error
	if err != nil {
		log.Error().Err(err).Interface("pedido", pedido).Msg("Failed to update pedido")
		return d.PedidoDetails(pedidoId)
	}

	return pedido, nil
}
