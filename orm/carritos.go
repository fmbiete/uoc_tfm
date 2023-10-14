package orm

import (
	"errors"

	"gorm.io/gorm"
)

func (d *Database) CarritoDelete(userId uint64) (Carrito, error) {
	var err error
	var carrito Carrito

	// Get the carrito id
	err = d.db.Select("id, usuario_id").Where("usuario_id = ?", userId).First(&carrito).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return carrito, err
	}

	if carrito.ID != 0 {
		// Delete lineas
		d.db.Unscoped().Model(&carrito).Association("CarritoLineas").Unscoped().Clear()
	}

	err = d.db.Save(&carrito).Error
	if err != nil {
		return carrito, err
	}
	return d.CarritoDetails(carrito.UsuarioID)
}

func (d *Database) CarritoDetails(userId uint64) (Carrito, error) {
	var carrito Carrito
	err := d.db.Preload("CarritoLineas").Where("usuario_id = ?", userId).First(&carrito).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No carrito, return one empty
		return Carrito{UsuarioID: userId}, nil
	}
	return carrito, err
}

func (d *Database) CarritoSave(carrito Carrito) (Carrito, error) {
	var err error

	// If there is no carrito ID, find the carrito for the user
	if carrito.ID == 0 {
		carrito, err = d.CarritoDetails(carrito.UsuarioID)
		if err != nil {
			return carrito, err
		}
	}

	// existing lineas: delete + insert
	lineas := carrito.CarritoLineas
	d.db.Unscoped().Model(&carrito).Association("CarritoLineas").Unscoped().Clear()
	carrito.CarritoLineas = lineas

	err = d.db.Save(&carrito).Error
	if err != nil {
		return carrito, err
	}
	return d.CarritoDetails(carrito.UsuarioID)
}
