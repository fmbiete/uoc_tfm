package orm

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

const errMsgReadConfig string = "Failed to read config"

func (d *Database) ConfigurationDetails() (Configuration, error) {
	var config Configuration
	err := d.db.First(&config).Error
	return config, err
}

func (d *Database) ConfigurationModify(config Configuration) (Configuration, error) {
	err := d.db.Model(&config).Updates(&config).Error
	if err != nil {
		return config, err
	}

	return d.ConfigurationDetails()
}

func (d *Database) configChangesAllowed() error {
	// Read ChangesTime
	var config Configuration
	err := d.db.Select("changes_time").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return err
	}

	// current time is before ChangesTime
	todayLimit := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.ChangesTime.Hour(), config.ChangesTime.Minute(), 0, 0, time.Now().Location())

	if time.Now().After(todayLimit) {
		return errors.New("kitchen is closed, no more orders or changes allowed")
	}

	return nil
}

func (d *Database) configSubvention() (float64, error) {
	// Read Subvention
	var config Configuration
	err := d.db.Select("subvention").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return 0, err
	}

	return config.Subvention, nil
}

func (d *Database) configTodayDelivery() (time.Time, error) {
	// Read DeliveryTime
	var config Configuration
	err := d.db.Select("delivery_time").First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg(errMsgReadConfig)
		return time.Now(), err
	}

	todayDelivery := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.DeliveryTime.Hour(), config.DeliveryTime.Minute(), 0, 0, time.Now().Location())
	return todayDelivery, nil
}
