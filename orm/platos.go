package orm

import (
	"errors"

	"gorm.io/gorm"
)

func (d *Database) PlatoCreate(plato Plato) (Plato, error) {
	err := d.db.Where("nombre = ?", plato.Nombre).First(&Plato{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := d.db.Create(&plato).Error
		return plato, err
	}

	if err != nil {
		return plato, err
	}

	// No error, we have found a matching plato - return duplicated error
	return plato, gorm.ErrDuplicatedKey
}

func (d *Database) PlatoDelete(platoId uint64) error {
	return d.db.Delete(&Plato{}, platoId).Error
}

func (d *Database) PlatoDetails(platoId uint64) (Plato, error) {
	var plato Plato
	err := d.db.Preload("Alergenos").Preload("Ingredientes").First(&plato, platoId).Error
	return plato, err
}

func (d *Database) PlatoList(usuarioId int64) ([]Plato, error) {
	var platos []Plato
	// TODO: filter by usuarioId if != -1, order by user sales
	// TODO: order by global sales
	err := d.db.Preload("Alergenos").Find(&platos).Error
	return platos, err
}

func (d *Database) PlatoModify(plato Plato) (Plato, error) {
	// replace alergenos and ingredientes - Update adds new records, but doesn't delete old ones
	alergenos := plato.Alergenos
	d.db.Unscoped().Model(&plato).Association("Alergenos").Unscoped().Clear()
	plato.Alergenos = alergenos

	ingredientes := plato.Ingredientes
	d.db.Unscoped().Model(&plato).Association("Ingredientes").Unscoped().Clear()
	plato.Ingredientes = ingredientes

	err := d.db.Updates(&plato).Error
	if err != nil {
		return plato, err
	}

	return d.PlatoDetails(uint64(plato.ID))
}

func (d *Database) platoCurrentPrecio(platoId uint64) (float64, error) {
	var err error
	var plato Plato

	err = d.db.Select("precio").First(&plato, platoId).Error
	return plato.Precio, err

	// TODO: promociones
}
