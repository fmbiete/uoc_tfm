package orm

import (
	"errors"
	"tfm_backend/models"

	"gorm.io/gorm"
)

func (d *Database) DishCreate(dish models.Dish) (models.Dish, error) {
	err := d.db.Where("name = ?", dish.Name).First(&models.Dish{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := d.db.Create(&dish).Error
		return dish, err
	}

	if err != nil {
		return dish, err
	}

	// No error, we have found a matching dish - return duplicated error
	return dish, gorm.ErrDuplicatedKey
}

func (d *Database) DishDelete(dishId uint64) error {
	return d.db.Delete(&models.Dish{}, dishId).Error
}

func (d *Database) DishDetails(dishId uint64) (models.Dish, error) {
	var dish models.Dish
	err := d.db.Preload("Allergens").Preload("Ingredients").First(&dish, dishId).Error
	return dish, err
}

func (d *Database) DishDislike(dishId uint64) error {
	return d.db.Exec(`UPDATE dishes SET dislike = dislike + 1 WHERE id = ?`, dishId).Error
}

func (d *Database) DishLike(dishId uint64) error {
	return d.db.Exec(`UPDATE dishes SET like = like + 1 WHERE id = ?`, dishId).Error
}

func (d *Database) DishList(userId int64, limit uint64, offset uint64) ([]models.Dish, error) {
	var dishes []models.Dish
	// TODO: filter by userId if != -1, order by user sales
	// TODO: order by global sales
	err := d.db.Preload("Allergens").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	return dishes, err
}

func (d *Database) DishModify(dish models.Dish) (models.Dish, error) {
	// replace alergenos and ingredientes - Update adds new records, but doesn't delete old ones
	alergenos := dish.Allergens
	d.db.Unscoped().Model(&dish).Association("Allergens").Unscoped().Clear()
	dish.Allergens = alergenos

	ingredientes := dish.Ingredients
	d.db.Unscoped().Model(&dish).Association("Ingredients").Unscoped().Clear()
	dish.Ingredients = ingredientes

	err := d.db.Updates(&dish).Error
	if err != nil {
		return dish, err
	}

	return d.DishDetails(uint64(dish.ID))
}

func (d *Database) dishCurrentCost(dishId uint64) (float64, error) {
	var err error
	var dish models.Dish

	err = d.db.Select("cost").First(&dish, dishId).Error
	return dish.Cost, err

	// TODO: promociones
}
