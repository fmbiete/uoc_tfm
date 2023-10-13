package orm

import (
	"errors"

	"gorm.io/gorm"
)

func (d *Database) PlatoCreate(plato Plato) (Plato, error) {
	err := d.db.Where("nombre = ?", plato.Nombre).First(&Plato{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := d.db.Create(&plato)
		return plato, result.Error
	}

	if err == nil {
		return plato, gorm.ErrDuplicatedKey
	}

	return plato, err
}

func (d *Database) PlatoDelete(platoId int64) error {
	return d.db.Delete(&Plato{}, platoId).Error
}

func (d *Database) PlatoDetails(platoId int64) (Plato, error) {
	var plato Plato
	result := d.db.Preload("Alergenos").Preload("Ingredientes").First(&plato, platoId)
	return plato, result.Error
}

func (d *Database) PlatoList(usuarioId int64) ([]Plato, error) {
	var platos []Plato
	// TODO: filter by usuarioId if != -1, order by user sales
	// TODO: order by global sales
	result := d.db.Preload("Alergenos").Find(&platos)
	return platos, result.Error
}

func (d *Database) PlatoModify(plato Plato) (Plato, error) {
	// replace alergenos and ingredientes - Update adds new records, but doesn't delete old ones
	alergenos := plato.Alergenos
	d.db.Unscoped().Model(&plato).Association("Alergenos").Unscoped().Clear()
	plato.Alergenos = alergenos

	ingredientes := plato.Ingredientes
	d.db.Unscoped().Model(&plato).Association("Ingredientes").Unscoped().Clear()
	plato.Ingredientes = ingredientes

	result := d.db.Updates(&plato)
	// returns only modified fields
	if result.Error == nil {
		return d.PlatoDetails(int64(plato.ID))
	}
	return plato, result.Error
}
