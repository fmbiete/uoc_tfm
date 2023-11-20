package orm

import (
	"errors"
	"fmt"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) DishCount() (int64, error) {
	var count int64
	err := d.db.Model(&models.Dish{}).Count(&count).Error
	return count, err
}

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
	err := d.db.Preload("Allergens", func(db *gorm.DB) *gorm.DB {
		return db.Order("allergens.name")
	}).Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredients.name")
	}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Order("categories.name")
	}).Preload("Promotions").
		First(&dish, dishId).Error
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

func (d *Database) DishFavourites(userId int64, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	if userId >= 0 {
		// Show my favourites
		err = d.db.Preload("Allergens", func(db *gorm.DB) *gorm.DB {
			return db.Order("allergens.name")
		}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Order("categories.name")
		}).Preload("Promotions").
			Joins("RIGHT JOIN dish_likes ON dish_likes.dish_id = dishes.id").Where(`dish_likes.user_id = ?`, userId).Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	}

	if len(dishes) == 0 {
		// Show global favourites
		err = d.db.Preload("Allergens", func(db *gorm.DB) *gorm.DB {
			return db.Order("allergens.name")
		}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Order("categories.name")
		}).Preload("Promotions").
			Order("likes desc").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	}

	return dishes, err
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

func (d *Database) DishList(searchTerm string, limit uint64, offset uint64) ([]models.Dish, error) {
	var err error
	var dishes []models.Dish

	scope := d.db
	if len(searchTerm) > 0 {
		scope = scope.Where("name ILIKE ?", fmt.Sprintf(`%%%s%%`, searchTerm))
	}
	err = scope.Preload("Promotions", func(db *gorm.DB) *gorm.DB {
		return db.Order("promotions.start_time DESC")
	}).Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredients.name")
	}).Preload("Allergens", func(db *gorm.DB) *gorm.DB {
		return db.Order("allergens.name")
	}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Order("categories.name")
	}).Order("name").Limit(int(limit)).Offset(int(offset)).Find(&dishes).Error
	if err != nil {
		log.Error().Err(err).Uint64("limit", limit).Uint64("offset", offset).Msg("Failed to list dishes")
		return dishes, err
	}

	return dishes, nil
}

func (d *Database) DishModify(dish models.Dish) (models.Dish, error) {
	var err error

	tx := d.db.Begin()
	defer tx.Rollback()

	// replace allergens and ingredients - Update adds new records, but doesn't delete old ones
	err = tx.Model(&dish).Association("Allergens").Replace(dish.Allergens)
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to replace dish allergens")
		return dish, err
	}

	err = tx.Model(&dish).Association("Ingredients").Replace(dish.Ingredients)
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to replace dish ingredients")
		return dish, err
	}

	err = tx.Model(&dish).Association("Categories").Replace(dish.Categories)
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to replace dish categories")
		return dish, err
	}

	err = tx.Updates(&dish).Error
	if err != nil {
		return dish, err
	}

	err = tx.Commit().Error
	if err != nil {
		log.Error().Err(err).Interface("dish", dish).Msg("Failed to commit modify dish")
	}

	return d.DishDetails(uint64(dish.ID))
}

func (d *Database) dishCurrentCost(dishId uint64) (float64, error) {
	var err error
	var dish models.Dish
	var promotion models.Promotion

	err = d.db.Select("cost").Where("dish_id = ? AND current_date BETWEEN start_time AND end_time", dishId).First(&promotion).Error
	if err == nil {
		// Dish has active Promotion
		return promotion.Cost, nil
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Dish doesn't have promotion
			err = d.db.Select("cost").First(&dish, dishId).Error
			return dish.Cost, err
		}
	}
	return 0, err
}
