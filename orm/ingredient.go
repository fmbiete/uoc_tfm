package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) IngredientCreate(ingredient models.Ingredient) (models.Ingredient, error) {
	err := d.db.Create(&ingredient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return d.IngredientFind(ingredient.Name)
		}

		log.Error().Err(err).Interface("ingredient", ingredient).Msg("Failed to create Ingredient")
		return models.Ingredient{}, err
	}

	return ingredient, nil
}

func (d *Database) IngredientDetails(ingredientId uint64) (models.Ingredient, error) {
	var ingredient models.Ingredient
	err := d.db.First(&ingredient, ingredientId).Error
	if err != nil {
		log.Error().Err(err).Uint64("ingredientId", ingredientId).Msg("Failed to detail Ingredient")
		return models.Ingredient{}, err
	}

	return ingredient, nil
}

func (d *Database) IngredientDelete(ingredientId uint64) error {
	var value uint64
	// Is the ingredient associated to any dish?
	res := d.db.Raw(`SELECT dish_id FROM dish_ingredients WHERE ingredient_id = ? LIMIT 1`, ingredientId).Scan(&value)
	if res.Error != nil {
		log.Error().Err(res.Error).Uint64("ingredientId", ingredientId).Msg("Failed to find Dish with Ingredient")
		return res.Error
	}

	if res.RowsAffected > 0 {
		log.Warn().Uint64("ingredientId", ingredientId).Msg("Dishes with Ingredient exist - we cannot remove it")
		return errors.New(`Dishes associated to this Ingredient exist - Remove the Dish association first`)
	}

	err := d.db.Unscoped().Where("id = ?", ingredientId).Delete(&models.Ingredient{}).Error
	if err != nil {
		log.Error().Err(err).Uint64("ingredientId", ingredientId).Msg("Failed to delete Ingredient")
		return err
	}

	return nil
}

func (d *Database) IngredientDishes(ingredientId uint64, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	err = d.db.Preload("Ingredients").Preload("Categories").Joins("RIGHT JOIN dish_categories ON dish_categories.dish_id = dishes.id").
		Where(`dish_categories.ingredient_id = ?`, ingredientId).
		Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	if err != nil {
		log.Error().Err(err).Uint64("ingredientId", ingredientId).Uint64("limit", limit).Uint64("offset", offset).Msg("Failed to read Ingredient Dishes")
		return dishes, err
	}

	return dishes, nil
}

func (d *Database) IngredientFind(name string) (models.Ingredient, error) {
	var ingredient models.Ingredient
	err := d.db.Where("name = ?", name).First(&ingredient).Error
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to find Ingredient")
		return models.Ingredient{}, err
	}

	return ingredient, nil
}

func (d *Database) IngredientList() ([]models.Ingredient, error) {
	var err error
	var categories []models.Ingredient

	err = d.db.Order("name").Find(&categories).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Ingredient")
		return categories, err
	}

	return categories, err
}

func (d *Database) IngredientModify(ingredient models.Ingredient) (models.Ingredient, error) {
	err := d.db.Save(&ingredient).Error
	if err != nil {
		log.Error().Err(err).Interface("ingredient", ingredient).Msg("Failed to update Ingredient")
		return models.Ingredient{}, err
	}

	return ingredient, nil
}
