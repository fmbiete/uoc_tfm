package orm

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

func (d *Database) PromocionCreate(promocion Promocion) (Promocion, error) {
	// don't allow 2 active promotions for the same plato
	err := d.db.Where("plato_id = ? and ? between inicio and fin", promocion.PlatoID, time.Now()).First(&Promocion{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result := d.db.Create(&promocion)
		return promocion, result.Error
	}

	if err == nil {
		return promocion, gorm.ErrDuplicatedKey
	}

	return promocion, err
}

func (d *Database) PromocionDelete(promocionId uint64) error {
	return d.db.Delete(&Promocion{}, promocionId).Error
}

func (d *Database) PromocionDetails(promocionId uint64) (Promocion, error) {
	var promocion Promocion
	result := d.db.First(&promocion, promocionId)
	return promocion, result.Error
}

func (d *Database) PromocionList() ([]Promocion, error) {
	var promociones []Promocion
	result := d.db.Where("? between inicio and fin", time.Now()).Find(&promociones)
	return promociones, result.Error
}

func (d *Database) PromocionModify(promocion Promocion) (Promocion, error) {
	result := d.db.Updates(&promocion)
	// returns only modified fields
	if result.Error == nil {
		return d.PromocionDetails(uint64(promocion.ID))
	}
	return promocion, result.Error
}
