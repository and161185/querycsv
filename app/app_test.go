package app

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {

	app := &app{}
	app.config = &config{}
	app.log = logrus.New()
	app.log.Formatter = new(logrus.JSONFormatter)

	dir, err := os.Getwd()
	if err != nil {
		app.log.Fatal(err)
	}
	app.log.Info(dir)

	configPath := "..\\config\\config.yaml"
	app.ReadConfigFile(configPath)

	assert.Equal(t, ".\\testing\\table.csv", app.config.TableFile)
	assert.Equal(t, 5, app.config.Timeout)
	assert.Equal(t, logrus.DebugLevel, app.log.Level)

	configPath = "..\\config\\config_err.yaml"
	app.ReadConfigFile(configPath)

	assert.Equal(t, "..\\testing\\table.csv", app.config.TableFile)
	assert.Equal(t, 5, app.config.Timeout)
	assert.Equal(t, logrus.InfoLevel, app.log.Level)
}
