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
		log.Error().Err(err).Interface("category", category).Msg("Failed to detail Category")
		return models.Category{}, err
	}

	return category, nil
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
