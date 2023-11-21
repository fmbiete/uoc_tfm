package orm

import (
	"errors"
	"fmt"
	"tfm_backend/models"
	"time"

	"github.com/rs/zerolog/log"
)

const errMsgReadConfig string = "Failed to read config"

func (d *Database) ConfigurationDetails() (models.Configuration, error) {
	var config models.Configuration
	err := d.db.First(&config).Error
	return config, err
}

func (d *Database) ConfigurationModify(config models.Configuration) (models.Configuration, error) {
	err := d.db.Model(&config).Updates(&config).Error
	if err != nil {
		return config, err
	}

	return d.ConfigurationDetails()
}

func (d *Database) configChangesAllowed(orderDelivery time.Time) error {
	if orderDelivery.Day() != time.Now().Day() ||
		orderDelivery.Month() != time.Now().Month() ||
		orderDelivery.Year() != time.Now().Year() {
		return errors.New("this order has been already delivered")
	}

	// Read ChangesTime
	var config models.Configuration
	err := d.db.Select("changes_time").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return err
	}

	// order time is before ChangesTime
	orderLimit := time.Date(orderDelivery.Year(), orderDelivery.Month(), orderDelivery.Day(),
		config.ChangesTime.Hour(), config.ChangesTime.Minute(), 0, 0, orderDelivery.Location())

	if time.Now().After(orderLimit) {
		return fmt.Errorf("kitchen is closed, no more orders or changes allowed after %s", orderLimit.Format("15:04"))
	}

	return nil
}

func (d *Database) configSubvention() (float64, error) {
	// Read Subvention
	var config models.Configuration
	err := d.db.Select("subvention").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return 0, err
	}

	return config.Subvention, nil
}

func (d *Database) configTodayDelivery() (time.Time, error) {
	// Read DeliveryTime
	var config models.Configuration
	err := d.db.Select("delivery_time").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return time.Now(), err
	}

	todayDelivery := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.DeliveryTime.Hour(), config.DeliveryTime.Minute(), 0, 0, time.Now().Location())
	return todayDelivery, nil
}
