package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) AllergenCreate(allergen models.Allergen) (models.Allergen, error) {
	err := d.db.Create(&allergen).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return d.AllergenFind(allergen.Name)
		}

		log.Error().Err(err).Interface("allergen", allergen).Msg("Failed to create Allergen")
		return models.Allergen{}, err
	}

	return allergen, nil
}

func (d *Database) AllergenDetails(allergenId uint64) (models.Allergen, error) {
	var allergen models.Allergen
	err := d.db.First(&allergen, allergenId).Error
	if err != nil {
		log.Error().Err(err).Uint64("allergenId", allergenId).Msg("Failed to detail Allergen")
		return models.Allergen{}, err
	}

	return allergen, nil
}

func (d *Database) AllergenDelete(allergenId uint64) error {
	var value uint64
	// Is the allergen associated to any dish?
	res := d.db.Raw(`SELECT dish_id FROM dish_allergens WHERE allergen_id = ? LIMIT 1`, allergenId).Scan(&value)
	if res.Error != nil {
		log.Error().Err(res.Error).Uint64("allergenId", allergenId).Msg("Failed to find Dish with Allergen")
		return res.Error
	}

	if res.RowsAffected > 0 {
		log.Warn().Uint64("allergenId", allergenId).Msg("Dishes with Allergen exist - we cannot remove it")
		return errors.New(`Dishes associated to this Allergen exist - Remove the Dish association first`)
	}

	err := d.db.Unscoped().Where("id = ?", allergenId).Delete(&models.Allergen{}).Error
	if err != nil {
		log.Error().Err(err).Uint64("allergenId", allergenId).Msg("Failed to delete Allergen")
		return err
	}

	return nil
}

func (d *Database) AllergenDishes(allergenId uint64, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	err = d.db.Preload("Allergens").Preload("Categories").Joins("RIGHT JOIN dish_categories ON dish_categories.dish_id = dishes.id").
		Where(`dish_categories.allergen_id = ?`, allergenId).
		Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	if err != nil {
		log.Error().Err(err).Uint64("allergenId", allergenId).Uint64("limit", limit).Uint64("offset", offset).Msg("Failed to read Allergen Dishes")
		return dishes, err
	}

	return dishes, nil
}

func (d *Database) AllergenFind(name string) (models.Allergen, error) {
	var allergen models.Allergen
	err := d.db.Where("name = ?", name).First(&allergen).Error
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to find Allergen")
		return models.Allergen{}, err
	}

	return allergen, nil
}

func (d *Database) AllergenList() ([]models.Allergen, error) {
	var err error
	var categories []models.Allergen

	err = d.db.Order("name").Find(&categories).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Allergen")
		return categories, err
	}

	return categories, err
}

func (d *Database) AllergenModify(allergen models.Allergen) (models.Allergen, error) {
	err := d.db.Save(&allergen).Error
	if err != nil {
		log.Error().Err(err).Interface("allergen", allergen).Msg("Failed to update Allergen")
		return models.Allergen{}, err
	}

	return allergen, nil
}
