package orm

import "tfm_backend/models"

func (d *Database) Login(username string, password string) (models.User, error) {
	var user models.User
	err := d.db.Where("email = ? AND password = ?", username, password).First(&user).Error
	// don't return the password hash
	user.Password = ""
	return user, err
}
