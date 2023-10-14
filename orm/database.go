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

	err = d.db.AutoMigrate(&Configuration{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Dish{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Promotion{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Order{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&OrderLine{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&Cart{})
	if err != nil {
		return err
	}
	err = d.db.AutoMigrate(&CartLine{})
	if err != nil {
		return err
	}

	return err
}
