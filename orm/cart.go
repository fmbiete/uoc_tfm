package orm

import (
	"errors"
	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (d *Database) CartDelete(userId uint64) (models.Cart, error) {
	var err error
	var cart models.Cart

	// Get the cart id
	err = d.db.Select("id, user_id").Where("user_id = ?", userId).First(&cart).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return cart, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart.UserID = userId
		err = d.db.Save(&cart).Error
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Msg("Failed to create new cart")
			return models.Cart{}, err
		}
	} else {
		// Delete lines
		err = d.db.Unscoped().Model(&cart).Association("CartLines").Unscoped().Clear()
		if err != nil {
			log.Error().Err(err).Uint64("userId", userId).Msg("Failed to delete cart lines from existing cart")
			return d.CartDetails(cart.UserID)
		}
	}

	return d.CartDetails(cart.UserID)
}

func (d *Database) CartDetails(userId uint64) (models.Cart, error) {
	var cart models.Cart
	err := d.db.Preload("CartLines").Where("user_id = ?", userId).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No cart, return one empty
		return models.Cart{UserID: userId}, nil
	}
	return cart, err
}

func (d *Database) CartSave(cart models.Cart) (models.Cart, error) {
	var err error

	// temp save of lines
	lines := cart.CartLines

	// If there is no cart ID, find the cart for the user
	if cart.ID == 0 {
		cart, err = d.CartDetails(cart.UserID)
		if err != nil {
			return cart, err
		}
	}

	// transaction block
	{
		tx := d.db.Begin()
		defer tx.Rollback()

		// existing lines: [hard] delete + insert
		err = tx.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartLine{}).Error
		if err != nil {
			log.Error().Err(err).Uint64("cartId", cart.ID).Msg("Failed to remove lines from cart")
		}
		cart.CartLines = lines

		err = tx.Save(&cart).Error
		if err != nil {
			return cart, err
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error().Err(err).Uint64("cartId", cart.ID).Msg("Failed to commit changes to cart")
			// Explicit rollback - next function is non-tx
			tx.Rollback()
			return d.CartDetails(cart.UserID)
		}
	}

	return d.CartDetails(cart.UserID)
}
