package orm

import (
	"errors"

	"gorm.io/gorm"
)

func (d *Database) CartDelete(userId uint64) (Cart, error) {
	var err error
	var cart Cart

	// Get the cart id
	err = d.db.Select("id, user_id").Where("user_id = ?", userId).First(&cart).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return cart, err
	}

	if cart.ID != 0 {
		// Delete lines
		d.db.Unscoped().Model(&cart).Association("CartLines").Unscoped().Clear()
	}

	err = d.db.Save(&cart).Error
	if err != nil {
		return cart, err
	}
	return d.CartDetails(cart.UserID)
}

func (d *Database) CartDetails(userId uint64) (Cart, error) {
	var cart Cart
	err := d.db.Preload("CartLines").Where("user_id = ?", userId).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No cart, return one empty
		return Cart{UserID: userId}, nil
	}
	return cart, err
}

func (d *Database) CartSave(cart Cart) (Cart, error) {
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

	// existing lines: delete + insert
	d.db.Unscoped().Model(&cart).Association("CartLines").Unscoped().Clear()
	cart.CartLines = lines

	err = d.db.Save(&cart).Error
	if err != nil {
		return cart, err
	}
	return d.CartDetails(cart.UserID)
}
