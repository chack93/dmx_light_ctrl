package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8080")
	viper.SetDefault("server.host", viper.GetString("host"))
	viper.SetDefault("server.port", viper.GetString("port"))

	viper.SetDefault("dmx.pin.a", 3)
	viper.SetDefault("dmx.pin.b", 2)
	viper.SetDefault("dmx.channel.total", 128)
	viper.SetDefault("dmx.protocol", "tm512ac")

	return nil
}
