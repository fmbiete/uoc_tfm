package orm

import (
	"fmt"
	"time"

	"tfm_backend/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	cfg        *models.ConfigDatabase
	siteAdmin  *models.User
	siteConfig *models.Configuration
	db         *gorm.DB
	models     []interface{}
}

func NewDatabase(cfg *models.Config) *Database {
	d := Database{cfg: &cfg.Database, siteAdmin: &cfg.SiteAdmin, siteConfig: &cfg.SiteConfig}

	d.models = append(d.models, &models.Configuration{})
	d.models = append(d.models, &models.User{})
	d.models = append(d.models, &models.Dish{})
	d.models = append(d.models, &models.Promotion{})
	d.models = append(d.models, &models.Order{})
	d.models = append(d.models, &models.OrderLine{})
	d.models = append(d.models, &models.DishDislike{})
	d.models = append(d.models, &models.DishLike{})

	return &d
}

func (d *Database) Setup() error {
	var err error

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%d sslmode=allow`,
		d.cfg.Host, d.cfg.User, d.cfg.Password, d.cfg.Database, d.cfg.Port)
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Error().Err(err).Msg("Failed to open gorm DB object")
		return err
	}

	sqlDb, err := d.db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to obtain SQL DB object")
		return err
	}

	log.Info().Msg("Configuring DB connection pool")
	sqlDb.SetConnMaxIdleTime(time.Hour)
	sqlDb.SetMaxIdleConns(d.cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(d.cfg.MaxOpenConns)

	if d.cfg.Reset {
		err = d.autoReset()
		if err != nil {
			return err
		}
	}

	err = d.autoMigrate()
	if err != nil {
		return err
	}

	if d.cfg.Reset {
		err = d.autoInit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) autoInit() error {
	var err error

	log.Info().Msg("Creating basic objects")

	err = d.db.Save(&d.siteAdmin).Error
	if err != nil {
		return err
	}

	err = d.db.Save(&d.siteConfig).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) autoReset() error {
	var err error

	log.Warn().Msg("Database Reset enabled - Destroying objects")

	for _, model := range d.models {
		err = d.db.Migrator().DropTable(model)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) autoMigrate() error {
	var err error

	log.Info().Msg("Automigrating ORM")

	for _, model := range d.models {
		err = d.db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}
