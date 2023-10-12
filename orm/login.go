package orm

func (d *Database) Login(username string, password string) (Usuario, error) {
	var usuario Usuario
	result := d.db.Where("email = ? AND password = ?", username, password).First(&usuario)
	return usuario, result.Error
}
