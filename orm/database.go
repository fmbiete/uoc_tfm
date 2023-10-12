package orm

import (
	"fmt"

	"tfm_backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	dsn string
	db  *gorm.DB
}

func NewDatabase(cfg config.ConfigDatabase) *Database {
	var d Database

	d.dsn = fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=allow`, cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	return &d
}

func (d *Database) Migrate() error {
	var err error
	d.db, err = gorm.Open(postgres.Open(d.dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Configuracion{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Usuario{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Plato{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Pedido{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&PedidoLinea{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Carrito{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&CarritoLinea{})
	if err != nil {
		return err
	}

	return err
}
