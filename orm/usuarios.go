package orm

import (
	"errors"

	"gorm.io/gorm"
)

func (d *Database) UsuarioCreate(user Usuario) (Usuario, error) {
	err := d.db.Where("email = ?", user.Email).First(&Usuario{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := d.db.Create(&user).Error
		return user, err
	}

	if err != nil {
		return user, err
	}

	// No error, we have found a matching user - return duplicated error
	return user, gorm.ErrDuplicatedKey
}

func (d *Database) UsuarioDelete(userId uint64) error {
	return d.db.Delete(&Usuario{}, userId).Error
}

func (d *Database) UsuarioDetails(userId uint64) (Usuario, error) {
	var user Usuario
	err := d.db.First(&user, userId).Error
	// Don't return the password hash
	user.Password = ""
	return user, err
}

func (d *Database) UsuarioModify(user Usuario) (Usuario, error) {
	err := d.db.Updates(&user).Error
	// Don't return the password hash
	user.Password = ""
	if err != nil {
		return user, err
	}

	return d.UsuarioDetails(uint64(user.ID))
}
