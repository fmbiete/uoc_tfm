package orm

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

func (d *Database) ConfiguracionDetails() (Configuracion, error) {
	var config Configuracion
	result := d.db.First(&config)
	return config, result.Error
}

func (d *Database) ConfiguracionModify(config Configuracion) (Configuracion, error) {
	result := d.db.Model(&config).Updates(&config)
	// returns only modified fields
	if result.Error == nil {
		return d.ConfiguracionDetails()
	}
	log.Error().Err(result.Error).Msg("HERE")
	return config, result.Error
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
		return errors.New("no more changes allowed")
	}

	return nil
}

func (d *Database) configPrecioSubvencion() (float32, error) {
	config, err := d.ConfiguracionDetails()
	if err != nil {
		return 0, err
	}

	return float32(config.PrecioSubvencion), nil
}

func (d *Database) configTodayEntrega() (time.Time, error) {
	config, err := d.ConfiguracionDetails()
	if err != nil {
		return time.Now(), err
	}

	todayEntrega := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), config.EntregaTime.Hour(), config.EntregaTime.Minute(), 0, 0, time.Now().Location())
	return todayEntrega, nil
}
