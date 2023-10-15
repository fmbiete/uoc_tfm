package orm

import (
	"errors"
	"tfm_backend/models"

	"gorm.io/gorm"
)

func (d *Database) UserCreate(user models.User) (models.User, error) {
	err := d.db.Where("email = ?", user.Email).First(&models.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = d.db.Create(&user).Error
		return user, err
	}

	if err != nil {
		return user, err
	}

	// No error, we have found a matching user - return duplicated error
	return user, gorm.ErrDuplicatedKey
}

func (d *Database) UserDelete(userId uint64) error {
	return d.db.Delete(&models.User{}, userId).Error
}

func (d *Database) UserDetails(userId uint64) (models.User, error) {
	var user models.User
	err := d.db.First(&user, userId).Error
	// Don't return the password hash
	user.Password = ""
	return user, err
}

func (d *Database) UserModify(user models.User) (models.User, error) {
	err := d.db.Updates(&user).Error
	// Don't return the password hash
	user.Password = ""
	if err != nil {
		return user, err
	}

	return d.UserDetails(uint64(user.ID))
}
