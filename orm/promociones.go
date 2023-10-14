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
		err := d.db.Create(&promocion).Error
		return promocion, err
	}

	if err != nil {
		return promocion, err
	}

	return promocion, gorm.ErrDuplicatedKey
}

func (d *Database) PromocionDelete(promocionId uint64) error {
	return d.db.Delete(&Promocion{}, promocionId).Error
}

func (d *Database) PromocionDetails(promocionId uint64) (Promocion, error) {
	var promocion Promocion
	err := d.db.First(&promocion, promocionId).Error
	return promocion, err
}

func (d *Database) PromocionList() ([]Promocion, error) {
	var promociones []Promocion
	err := d.db.Where("? between inicio and fin", time.Now()).Find(&promociones).Error
	return promociones, err
}

func (d *Database) PromocionModify(promocion Promocion) (Promocion, error) {
	err := d.db.Updates(&promocion).Error
	if err != nil {
		return promocion, err
	}

	return d.PromocionDetails(uint64(promocion.ID))
}
