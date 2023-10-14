package orm

func (d *Database) Login(username string, password string) (User, error) {
	var user User
	err := d.db.Where("email = ? AND password = ?", username, password).First(&user).Error
	// don't return the password hash
	user.Password = ""
	return user, err
}
