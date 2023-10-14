package orm

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

func (d *Database) ConfiguracionDetails() (Configuracion, error) {
	var config Configuracion
	err := d.db.First(&config).Error
	return config, err
}

func (d *Database) ConfiguracionModify(config Configuracion) (Configuracion, error) {
	err := d.db.Model(&config).Updates(&config).Error
	if err != nil {
		return config, err
	}

	return d.ConfiguracionDetails()
}

func (d *Database) configChangesAllowed() error {
	config, err := d.ConfiguracionDetails()
	if err != nil {
		return err
	}

	// current time is before cambiosTime
	todayLimit := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.CambiosTime.Hour(), config.CambiosTime.Minute(), 0, 0, time.Now().Location())

	if time.Now().After(todayLimit) {
		log.Info().Msg("CambiosTime has passed for today")
		return errors.New("kitchen is closed, no more pedidos or changes allowed")
	}

	return nil
}

func (d *Database) configPrecioSubvencion() (float64, error) {
	config, err := d.ConfiguracionDetails()
	if err != nil {
		return 0, err
	}

	return config.PrecioSubvencion, nil
}

func (d *Database) configTodayEntrega() (time.Time, error) {
	config, err := d.ConfiguracionDetails()
	if err != nil {
		return time.Now(), err
	}

	todayEntrega := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.EntregaTime.Hour(), config.EntregaTime.Minute(), 0, 0, time.Now().Location())
	return todayEntrega, nil
}
