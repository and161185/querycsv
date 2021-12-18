package app

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"querycsv/csvreader"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type csvreaderI interface {
	FindRows(context.Context) []map[string]string
}

type config struct {
	TableFile string `yaml:"tableFile"`
	Timeout   int    `yaml:"timeout"`
	LogLevel  string `yaml:"logLevel"`
}

type app struct {
	log    *logrus.Logger
	config *config
	query  string
	done   chan bool
}

const defaultLogLevel = logrus.InfoLevel
const defaultTimeout = 5

func NewApp(log *logrus.Logger) *app {

	log.Info("Initializing app")

	a := &app{
		log:    log,
		config: &config{},
		done:   make(chan bool),
	}

	a.log.Info("app initialized")
	return a
}

func (app *app) Done() chan bool {
	return app.done
}

func (app *app) ReadConfigFile(configPath string) {

	app.log.Info("loading settings")

	var logLevel logrus.Level

	_, err := os.Stat(configPath)

	defer func() {
		if err != nil {

			app.log.Errorf("Couldn't read config file %s , got %v", configPath, err)
			app.log.Info("Use app's default settings")

			logLevel = defaultLogLevel
			app.log.Infof("logLevel's default value %v is setted", defaultLogLevel)

			app.config.Timeout = defaultTimeout
			app.log.Infof("Timout's default value %v is setted", defaultTimeout)
		}

		if app.config.Timeout == 0 {
			app.config.Timeout = defaultTimeout
			app.log.Infof("Timeout can't be 0. Default value %v is setted", defaultTimeout)
		}

		app.log.Level = logLevel
		app.log.Info("Settings loaded")

	}()

	if err != nil {
		return
	}

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, app.config)
	if err != nil {
		return
	}

	app.log.Infof("string logrus level: %s", app.config.LogLevel)
	level, err := logrus.ParseLevel(app.config.LogLevel)
	if err != nil {
		return
	}
	logLevel = level

}

func (app *app) GetConfig() config {
	return (*app.config)
}

func (app *app) GetQueryString() {
loop:
	for {

		fmt.Println("Enter query:")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			app.query = scanner.Text()
			//app.query = "salary < 10000 or name = \"Ann\""

			var re = regexp.MustCompile(`[^=]=\s*"`)
			app.query = re.ReplaceAllString(app.query, ` == "`)
			re = regexp.MustCompile(`"\s*=[^=]`)
			app.query = re.ReplaceAllString(app.query, `" == `)
			break loop
		}

	}
}
func (app *app) Timeout() time.Duration {
	return time.Second * time.Duration(app.config.Timeout)
}

func (app *app) Run(ctx context.Context) {
	app.log.Info("guery: ", app.query)

	var reader csvreaderI
	reader = csvreader.NewReader(app.config.TableFile, app.log, app.query)

	rows := reader.FindRows(ctx)

	if len(rows) == 0 {
		fmt.Println("no results fo query ", app.query)
		return
	}

	ShowTable(rows)

	app.done <- true
}

func ShowTable(rows []map[string]string) {

	table := tablewriter.NewWriter(os.Stdout)
	var header []string

	for k := range rows[0] {
		header = append(header, k)
	}
	table.SetHeader(header)

	for i := 0; i < len(rows); i++ {
		row := rows[i]

		var slice []string
		for h := 0; h < len(header); h++ {
			slice = append(slice, row[header[h]])
		}
		table.Append(slice)
	}
	table.Render() // Send output
}
