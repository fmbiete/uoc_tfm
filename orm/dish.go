package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
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

func (d *Database) DishDislike(userId uint64, dishId uint64) error {
	var err error
	tx := d.db.Begin()
	defer tx.Rollback()

	err = tx.Create(&models.DishDislike{UserID: userId, DishID: dishId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Dislike already exists")
			return nil
		}

		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to create dislike")
		return err
	}

	// If the user liked this dish before, remove it
	result := tx.Unscoped().Where(`user_id = ? AND dish_id = ?`, userId, dishId).Delete(&models.DishLike{})
	if result.Error != nil {
		log.Error().Err(result.Error).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to destroy like")
		return result.Error
	}
	if result.RowsAffected > 0 {
		err = tx.Exec(`UPDATE dishes SET likes = likes - 1 WHERE id = ?`, dishId).Error
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to decr like count")
			return err
		}
	}

	err = tx.Exec(`UPDATE dishes SET dislikes = dislikes + 1 WHERE id = ?`, dishId).Error
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to incr dislike count")
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to commit tx dislike")
		return err
	}

	return nil
}

func (d *Database) DishLike(userId uint64, dishId uint64) error {
	var err error
	tx := d.db.Begin()
	defer tx.Rollback()

	err = tx.Create(&models.DishLike{UserID: userId, DishID: dishId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			log.Info().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Like already exists")
			return nil
		}

		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to create like")
		return err
	}

	// If the user disliked this dish before, remove it
	result := tx.Unscoped().Where(`user_id = ? AND dish_id = ?`, userId, dishId).Delete(&models.DishDislike{})
	if result.Error != nil {
		log.Error().Err(result.Error).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to destroy dislike")
		return result.Error
	}
	if result.RowsAffected > 0 {
		err = tx.Exec(`UPDATE dishes SET dislikes = dislikes - 1 WHERE id = ?`, dishId).Error
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to decr dislike count")
			return err
		}
	}

	err = tx.Exec(`UPDATE dishes SET likes = likes + 1 WHERE id = ?`, dishId).Error
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to incr like count")
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error().Err(err).Uint64("userId", userId).Uint64("dishId", dishId).Msg("Failed to commit tx like")
		return err
	}

	return nil
}

func (d *Database) DishList(userId int64, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	if userId >= 0 {
		// Show my favourites
		err = d.db.Preload("Allergens").Joins("RIGHT JOIN dish_likes ON dish_likes.dish_id = dishes.id").Where(`dish_likes.user_id = ?`, userId).Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	}

	if len(dishes) == 0 {
		// Show global favourites
		err = d.db.Preload("Allergens").Order("likes desc").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	}

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
