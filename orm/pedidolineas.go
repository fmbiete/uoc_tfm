package orm

import "github.com/rs/zerolog/log"

func (d *Database) PedidoLineaCreate(pedidoId uint64, lineaCarrito CarritoLinea) (Pedido, error) {
	pedido, err := d.PedidoDetails(pedidoId)
	if err != nil {
		log.Error().Err(err).Uint64("pedidoId", pedidoId).Msg("Failed to read pedido")
		return pedido, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return pedido, err
	}

	// Convert linea carrito to linea pedido
	var linea PedidoLinea = PedidoLinea{Cantidad: lineaCarrito.Cantidad, PedidoID: pedidoId}

	var plato Plato
	err = d.db.Select("nombre").First(&plato, lineaCarrito.PlatoID).Error
	if err != nil {
		log.Error().Err(err).Interface("linea", linea).Msg("Failed to read plato from carrito linea")
		return pedido, err
	}
	linea.Nombre = plato.Nombre

	var precioUnidad float64
	precioUnidad, err = d.platoCurrentPrecio(lineaCarrito.PlatoID)
	if err != nil {
		log.Error().Err(err).Uint64("platoId", lineaCarrito.PlatoID).Msg("Failed to read precio from plato")
		return pedido, err
	}
	linea.PrecioUnidad = precioUnidad

	err = d.db.Save(&linea).Error
	if err != nil {
		log.Error().Err(err).Interface("linea", linea).Msg("Failed to save linea pedido")
		return pedido, err
	}

	return d.pedidoUpdatePrecio(pedidoId)
}

func (d *Database) PedidoLineaDelete(pedidoId uint64, lineaId uint64) (Pedido, error) {
	var err error
	var pedido Pedido

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return pedido, err
	}

	err = d.db.Delete(&PedidoLinea{}, lineaId).Error
	if err != nil {
		log.Error().Err(err).Uint64("lineaId", lineaId).Msg("Failed to delete linea pedido")
		return pedido, err
	}

	return d.pedidoUpdatePrecio(pedidoId)
}
func (d *Database) PedidoLineaModify(linea PedidoLinea) (Pedido, error) {
	pedido, err := d.PedidoDetails(linea.PedidoID)
	if err != nil {
		log.Error().Err(err).Uint64("pedidoId", linea.PedidoID).Msg("Failed to read pedido")
		return pedido, err
	}

	// Changes must be done before kitchen starts preparing the food
	err = d.configChangesAllowed()
	if err != nil {
		return pedido, err
	}

	// existing linea - update cantidad ONLY
	err = d.db.Model(&linea).Update("cantidad", linea.Cantidad).Error
	if err != nil {
		log.Error().Err(err).Interface("linea", linea).Msg("Failed to save pedido linea")
		return pedido, err
	}

	// use save to update or create line as required
	return d.pedidoUpdatePrecio(linea.PedidoID)
}
