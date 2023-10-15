package orm

import (
	"errors"
	"tfm_backend/models"
	"time"

	"gorm.io/gorm"
)

func (d *Database) PromotionCreate(promotion models.Promotion) (models.Promotion, error) {
	// don't allow 2 active promotions for the same dish
	err := d.db.Where("dish_id = ? and ? between start_time and end_time", promotion.DishID, time.Now()).First(&models.Promotion{}).Error
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
	return d.db.Delete(&models.Promotion{}, promotionId).Error
}

func (d *Database) PromotionDetails(promotionId uint64) (models.Promotion, error) {
	var promotion models.Promotion
	err := d.db.First(&promotion, promotionId).Error
	return promotion, err
}

func (d *Database) PromotionList() ([]models.Promotion, error) {
	var promotiones []models.Promotion
	err := d.db.Where("? between start_time and end_time", time.Now()).Find(&promotiones).Error
	return promotiones, err
}

func (d *Database) PromotionModify(promotion models.Promotion) (models.Promotion, error) {
	err := d.db.Updates(&promotion).Error
	if err != nil {
		return promotion, err
	}

	return d.PromotionDetails(uint64(promotion.ID))
}
