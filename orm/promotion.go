package orm

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

func (d *Database) PromotionCreate(promotion Promotion) (Promotion, error) {
	// don't allow 2 active promotions for the same dish
	err := d.db.Where("dish_id = ? and ? between start and end", promotion.DishID, time.Now()).First(&Promotion{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := d.db.Create(&promotion).Error
		return promotion, err
	}

	if err != nil {
		return promotion, err
	}

	return promotion, gorm.ErrDuplicatedKey
}

func (d *Database) PromotionDelete(promotionId uint64) error {
	return d.db.Delete(&Promotion{}, promotionId).Error
}

func (d *Database) PromotionDetails(promotionId uint64) (Promotion, error) {
	var promotion Promotion
	err := d.db.First(&promotion, promotionId).Error
	return promotion, err
}

func (d *Database) PromotionList() ([]Promotion, error) {
	var promotiones []Promotion
	err := d.db.Where("? between start and end", time.Now()).Find(&promotiones).Error
	return promotiones, err
}

func (d *Database) PromotionModify(promotion Promotion) (Promotion, error) {
	err := d.db.Updates(&promotion).Error
	if err != nil {
		return promotion, err
	}

	return d.PromotionDetails(uint64(promotion.ID))
}
