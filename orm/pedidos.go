package orm

import (
	"gorm.io/gorm"
)

func (d *Database) PedidoCreate(pedido Pedido) (Pedido, error) {
	// Pedidos must be done before kitchen starts preparing the food
	err := d.configChangesAllowed()
	if err != nil {
		return pedido, err
	}

	// Calculate entrega
	pedido.Entrega, err = d.configTodayEntrega()
	if err != nil {
		return pedido, err
	}

	pedido, err = d.pedidoCalcularPrecio(pedido)
	if err != nil {
		return pedido, err
	}

	result := d.db.Create(&pedido)
	return pedido, result.Error
}

func (d *Database) PedidoDelete(pedidoId uint64) error {
	return d.db.Delete(&Pedido{}, pedidoId).Error
}

func (d *Database) PedidoDetails(pedidoId uint64) (Pedido, error) {
	var pedido Pedido
	result := d.db.Preload("PedidoLineas").First(&pedido, pedidoId)
	return pedido, result.Error
}

func (d *Database) PedidoList(usuarioId int64) ([]Pedido, error) {
	var pedidos []Pedido
	var result *gorm.DB
	// soft-deleted pedidos are cancelled
	if usuarioId == -1 {
		result = d.db.Unscoped().Find(&pedidos)
	} else {
		result = d.db.Unscoped().Where("usuario_id = ?", usuarioId).Find(&pedidos)
	}
	return pedidos, result.Error
}

func (d *Database) PedidoModify(pedido Pedido) (Pedido, error) {
	// Pedidos must be done before kitchen starts preparing the food
	err := d.configChangesAllowed()
	if err != nil {
		return pedido, err
	}

	pedido, err = d.pedidoCalcularPrecio(pedido)
	if err != nil {
		return pedido, err
	}

	// replace pedido lineas - Update adds new records, but doesn't delete old ones
	lineas := pedido.PedidoLineas
	d.db.Unscoped().Model(&pedido).Association("PedidoLineas").Unscoped().Clear()
	pedido.PedidoLineas = lineas

	result := d.db.Updates(&pedido)
	// returns only modified fields
	if result.Error == nil {
		return d.PedidoDetails(uint64(pedido.ID))
	}
	return pedido, result.Error
}

func (d *Database) pedidoCalcularPrecio(pedido Pedido) (Pedido, error) {
	// Allowed: multiple pedidos per day
	// TODO: apply subvencion only in the first pedido
	subvencion, err := d.configPrecioSubvencion()
	if err != nil {
		return pedido, err
	}

	// Calculate total
	for i, _ := range pedido.PedidoLineas {
		pedido.PrecioTotal += float32(pedido.PedidoLineas[i].Cantidad) * pedido.PedidoLineas[i].PrecioUnidad
	}
	// Apply subvencion
	pedido.PrecioPagar = pedido.PrecioTotal - subvencion
	if pedido.PrecioPagar < 0 {
		pedido.PrecioPagar = 0
	}

	return pedido, nil
}
