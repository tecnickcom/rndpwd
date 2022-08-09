package cli

import (
	"testing"

	"github.com/nexmoinc/gosrvlib/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tecnickcom/rndpwd/internal/validator"
)

func Test_appConfig_SetDefaults(t *testing.T) {
	t.Parallel()

	v := viper.New()
	c := &appConfig{}
	c.SetDefaults(v)

	require.True(t, v.GetBool("enabled"))
	require.NotEmpty(t, v.GetString("monitoring_address"))
	require.NotEmpty(t, v.GetString("public_address"))
	require.Equal(t, 8, len(v.AllKeys()))
}

func getValidTestConfig() appConfig {
	return appConfig{
		BaseConfig: config.BaseConfig{
			Log: config.LogConfig{
				Level:   "DEBUG",
				Format:  "CONSOLE",
				Network: "tcp",
				Address: "127.0.0.1:1234",
			},
		},
		Enabled:           true,
		MonitoringAddress: ":1233",
		PublicAddress:     ":1231",
		Ipify: ipifyConfig{
			Address: "https://test.ipify.url.invalid",
			Timeout: 1,
		},
		Random: randomConfig{
			Charset:  validator.ValidCharset,
			Length:   16,
			Quantity: 3,
		},
	}
}

func Test_appConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fcfg    func(cfg appConfig) appConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			fcfg:    func(cfg appConfig) appConfig { return cfg },
			wantErr: false,
		},
		{
			name:    "empty log.level",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Level = ""; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid log.level",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Level = "WRONG_LOG_LEVEL"; return cfg },
			wantErr: true,
		},
		{
			name:    "empty log.format",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Format = ""; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid log.format",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Format = "WRONG_LOG_FORMAT"; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid log.network",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Network = "WRONG_LOG_NETWORK"; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid log.address",
			fcfg:    func(cfg appConfig) appConfig { cfg.Log.Address = "-WRONG_LOG_ADDRESS-"; return cfg },
			wantErr: true,
		},
		{
			name:    "empty monitoring_address",
			fcfg:    func(cfg appConfig) appConfig { cfg.MonitoringAddress = ""; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid monitoring_address",
			fcfg:    func(cfg appConfig) appConfig { cfg.MonitoringAddress = "-WRONG_MONITORING_ADDRESS-"; return cfg },
			wantErr: true,
		},
		{
			name:    "empty public_address",
			fcfg:    func(cfg appConfig) appConfig { cfg.PublicAddress = ""; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid public_address",
			fcfg:    func(cfg appConfig) appConfig { cfg.PublicAddress = "-WRONG_PUBLIC_ADDRESS-"; return cfg },
			wantErr: true,
		},
		{
			name:    "empty ipify.address",
			fcfg:    func(cfg appConfig) appConfig { cfg.Ipify.Address = ""; return cfg },
			wantErr: true,
		},
		{
			name:    "invalid ipify.address",
			fcfg:    func(cfg appConfig) appConfig { cfg.Ipify.Address = "-WRONG_IPIFY_ADDRESS-"; return cfg },
			wantErr: true,
		},
		{
			name:    "empty ipify.timeout",
			fcfg:    func(cfg appConfig) appConfig { cfg.Ipify.Timeout = 0; return cfg },
			wantErr: true,
		},

		{
			name:    "too small random.charset",
			fcfg:    func(cfg appConfig) appConfig { cfg.Random.Charset = "a"; return cfg },
			wantErr: true,
		},
		{
			name: "too big random.charset",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Random.Charset = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
				return cfg
			},
			wantErr: true,
		},
		{
			name: "invalid random.charset",
			fcfg: func(cfg appConfig) appConfig {
				cfg.Random.Charset = " \xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
				return cfg
			},
			wantErr: true,
		},
		{
			name:    "too small random.length",
			fcfg:    func(cfg appConfig) appConfig { cfg.Random.Length = 1; return cfg },
			wantErr: true,
		},
		{
			name:    "empty random.quantity",
			fcfg:    func(cfg appConfig) appConfig { cfg.Random.Quantity = 0; return cfg },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := tt.fcfg(getValidTestConfig())
			if err := cfg.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
