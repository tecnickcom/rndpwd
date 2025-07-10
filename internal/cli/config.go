package cli

import (
	"github.com/Vonage/gosrvlib/pkg/config"
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
	appLongDesc = "Web-Service Random Password Generator"

	// fieldTagName is the name of the tag containing the original JSON field name.
	fieldTagName = "mapstructure"
)

type cfgServer struct {
	Address string `mapstructure:"address" validate:"required,hostname_port"`
	Timeout int    `mapstructure:"timeout" validate:"required,min=1"`
}

type cfgServerMonitoring cfgServer

type cfgServerPublic cfgServer

// cfgServers contains the configuration for all exposed servers.
type cfgServers struct {
	Monitoring cfgServerMonitoring `mapstructure:"monitoring" validate:"required"`
	Public     cfgServerPublic     `mapstructure:"public"     validate:"required"`
}

type cfgClientIpify struct {
	Address string `mapstructure:"address" validate:"required,url"`
	Timeout int    `mapstructure:"timeout" validate:"required,min=1"`
}

// cfgClients contains the configuration for all external clients.
type cfgClients struct {
	Ipify cfgClientIpify `mapstructure:"ipify" validate:"required"`
}

// randomConfig contains the random generator configuration.
type randomConfig struct {
	Charset  string `mapstructure:"charset"  validate:"required,min=1,max=256,rndcharset"`
	Length   int    `mapstructure:"length"   validate:"required,min=1,max=4096"`
	Quantity int    `mapstructure:"quantity" validate:"required,min=1,max=1000"`
}

// appConfig contains the full application configuration.
type appConfig struct {
	config.BaseConfig `mapstructure:",squash" validate:"required"`

	Enabled bool         `mapstructure:"enabled"`
	Servers cfgServers   `mapstructure:"servers" validate:"required"`
	Clients cfgClients   `mapstructure:"clients" validate:"required"`
	Random  randomConfig `mapstructure:"random"  validate:"required"`
}

// SetDefaults sets the default configuration values in Viper.
func (c *appConfig) SetDefaults(v config.Viper) {
	v.SetDefault("enabled", true)

	v.SetDefault("servers.monitoring.address", ":8072")
	v.SetDefault("servers.monitoring.timeout", 60)

	v.SetDefault("servers.public.address", ":8071")
	v.SetDefault("servers.public.timeout", 60)

	v.SetDefault("clients.ipify.address", "https://api.ipify.org")
	v.SetDefault("clients.ipify.timeout", 1)

	v.SetDefault("random.charset", validator.ValidCharset)
	v.SetDefault("random.length", 32)
	v.SetDefault("random.quantity", 10)
}

// Validate performs the validation of the configuration values.
func (c *appConfig) Validate() error {
	v, _ := validator.New(fieldTagName)

	return v.ValidateStruct(c) //nolint:wrapcheck
}
