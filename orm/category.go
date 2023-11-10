package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) CategoryCreate(category models.Category) (models.Category, error) {
	err := d.db.Create(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return d.CategoryFind(category.Name)
		}

		log.Error().Err(err).Interface("category", category).Msg("Failed to create Category")
		return models.Category{}, err
	}

	return category, nil
}

func (d *Database) CategoryDetails(categoryId uint64) (models.Category, error) {
	var category models.Category
	err := d.db.First(&category, categoryId).Error
	if err != nil {
		log.Error().Err(err).Uint64("categoryId", categoryId).Msg("Failed to detail Category")
		return models.Category{}, err
	}

	return category, nil
}

func (d *Database) CategoryDelete(categoryId uint64) error {
	var value uint64
	// Is the category associated to any dish?
	res := d.db.Raw(`SELECT dish_id FROM dish_categories WHERE category_id = ? LIMIT 1`, categoryId).Scan(&value)
	if res.Error != nil {
		log.Error().Err(res.Error).Uint64("categoryId", categoryId).Msg("Failed to find Dish with Category")
		return res.Error
	}

	if res.RowsAffected > 0 {
		log.Warn().Uint64("categoryId", categoryId).Msg("Dishes with Category exist - we cannot remove it")
		return errors.New(`Dishes associated to this Category exist - Remove the Dish association first`)
	}

	err := d.db.Unscoped().Where("id = ?", categoryId).Delete(&models.Category{}).Error
	if err != nil {
		log.Error().Err(err).Uint64("categoryId", categoryId).Msg("Failed to delete Category")
		return err
	}

	return nil
}

func (d *Database) CategoryDishes(categoryId uint64, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	err = d.db.Preload("Allergens").Preload("Categories").Joins("RIGHT JOIN dish_categories ON dish_categories.dish_id = dishes.id").
		Where(`dish_categories.category_id = ?`, categoryId).
		Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	if err != nil {
		log.Error().Err(err).Uint64("categoryId", categoryId).Uint64("limit", limit).Uint64("offset", offset).Msg("Failed to read Category Dishes")
		return dishes, err
	}

	return dishes, nil
}

func (d *Database) CategoryFind(name string) (models.Category, error) {
	var category models.Category
	err := d.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to find Category")
		return models.Category{}, err
	}

	return category, nil
}

func (d *Database) CategoryList() ([]models.Category, error) {
	var err error
	var categories []models.Category

	err = d.db.Order("name").Find(&categories).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Category")
		return categories, err
	}

	return categories, err
}

func (d *Database) CategoryModify(category models.Category) (models.Category, error) {
	err := d.db.Save(&category).Error
	if err != nil {
		log.Error().Err(err).Interface("category", category).Msg("Failed to update Category")
		return models.Category{}, err
	}

	return category, nil
}
