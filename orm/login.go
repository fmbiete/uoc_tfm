package orm

func (d *Database) Login(username string, password string) (Usuario, error) {
	var usuario Usuario
	err := d.db.Where("email = ? AND password = ?", username, password).First(&usuario).Error
	// don't return the password hash
	usuario.Password = ""
	return usuario, err
}
