package cli

import (
	"github.com/nexmoinc/gosrvlib/pkg/config"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

const (
	// AppName is the name of the application executable.
	AppName = "rndpwd"

	// appEnvPrefix is the prefix of the configuration environment variables.
	appEnvPrefix = "RNDPWD"

	// appShortDesc is the short description of the application.
	appShortDesc = "Web-Service Random Password Generator"

	// appLongDesc is the long description of the application.
	appLongDesc = "Web-Service Random Password Generator - gosrvlib example"

	// fieldTagName is the name of the tag containing the original JSON field name.
	fieldTagName = "mapstructure"
)

// ipifyConfig contains ipify client configuration.
type ipifyConfig struct {
	Address string `mapstructure:"address" validate:"required,url"`
	Timeout int    `mapstructure:"timeout" validate:"required,min=1"`
}

// randomConfig contains the random generator configuration.
type randomConfig struct {
	Charset  string `mapstructure:"charset" validate:"required,min=2,max=92,rndcharset"`
	Length   int    `mapstructure:"length" validate:"required,min=2,max=64"`
	Quantity int    `mapstructure:"quantity" validate:"required,min=1,max=10"`
}

// appConfig contains the full application configuration.
type appConfig struct {
	config.BaseConfig `mapstructure:",squash" validate:"required"`
	Enabled           bool         `mapstructure:"enabled"`
	MonitoringAddress string       `mapstructure:"monitoring_address" validate:"required,hostname_port"`
	PublicAddress     string       `mapstructure:"public_address" validate:"required,hostname_port"`
	Ipify             ipifyConfig  `mapstructure:"ipify" validate:"required"`
	Random            randomConfig `mapstructure:"random" validate:"required"`
}

// SetDefaults sets the default configuration values in Viper.
func (c *appConfig) SetDefaults(v config.Viper) {
	v.SetDefault("enabled", true)

	// Setting the default monitoring_address port to the same as service_port will start a single HTTP server
	v.SetDefault("monitoring_address", ":8072")
	v.SetDefault("public_address", ":8071")

	v.SetDefault("ipify.address", "https://api.ipify.org")
	v.SetDefault("ipify.timeout", 1)

	v.SetDefault("random.charset", validator.ValidCharset)
	v.SetDefault("random.length", 32)
	v.SetDefault("random.quantity", 10)
}

// Validate performs the validation of the configuration values.
func (c *appConfig) Validate() error {
	v, _ := validator.New(fieldTagName)

	return v.ValidateStruct(c) //nolint:wrapcheck
}
