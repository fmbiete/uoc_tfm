package orm

import (
	"errors"
	"tfm_backend/models"
	"time"

	"gorm.io/gorm"
)

func (d *Database) PromotionCount(activeOnly bool) (int64, error) {
	var count int64
	scope := d.db.Model(&models.Promotion{})
	if activeOnly {
		scope = scope.Where("? BETWEEN start_time AND end_time", time.Now())
	}
	err := scope.Count(&count).Error
	return count, err
}

func (d *Database) PromotionCreate(promotion models.Promotion) (models.Promotion, error) {
	// don't allow 2 overlapping promotions for the same dish
	err := d.db.Where("dish_id = ? and (? between start_time and end_time or ? between start_time and end_time)",
		promotion.DishID, promotion.StartTime, promotion.EndTime).
		First(&models.Promotion{}).Error
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

func (d *Database) PromotionList(activeOnly bool, limit uint64, offset uint64) ([]models.Promotion, error) {
	var promotions []models.Promotion
	scope := d.db
	if activeOnly {
		scope = scope.Where("? BETWEEN start_time AND end_time", time.Now())
	}
	err := scope.Preload("Dish").Joins("LEFT JOIN dishes ON promotions.dish_id = dishes.id").
		Preload("Dish.Allergens").
		Order("start_time DESC").Limit(int(limit)).Offset(int(offset)).Find(&promotions).Error
	return promotions, err
}

func (d *Database) PromotionModify(promotion models.Promotion) (models.Promotion, error) {
	err := d.db.Updates(&promotion).Error
	if err != nil {
		return promotion, err
	}

	return d.PromotionDetails(uint64(promotion.ID))
}
