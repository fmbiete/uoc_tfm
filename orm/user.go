package orm

import (
	"errors"
	"fmt"
	"tfm_backend/models"

	"gorm.io/gorm"
)

func (d *Database) UserCount(includeAdmin bool) (int64, error) {
	var count int64
	scope := d.db.Model(&models.User{})
	if !includeAdmin {
		scope = scope.Where("is_admin = false")
	}
	err := scope.Count(&count).Error
	return count, err
}

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

func (d *Database) UserList(searchTerm string, limit uint64, offset uint64) ([]models.User, error) {
	var users []models.User

	scope := d.db
	if len(searchTerm) > 0 {
		filter := fmt.Sprintf(`%%%s%%`, searchTerm)
		scope = scope.Where("name ILIKE ? OR surname ILIKE ? OR email ILIKE ?", filter, filter, filter)
	}

	err := scope.Order("is_admin DESC, name, email").Limit(int(limit)).Offset(int(offset)).Find(&users).Error
	return users, err
}

func (d *Database) UserModify(user models.User) (models.User, error) {
	err := d.db.Updates(&user).Error
	// Don't return the password hash
	user.Password = ""
	if err != nil {
		return user, err
	}

	if !user.IsAdmin {
		// Update admin flag - gorm will not update false
		err = d.db.Model(&user).Updates(map[string]interface{}{"is_admin": false}).Error
		if err != nil {
			return user, err
		}

	}

	return d.UserDetails(uint64(user.ID))
}
