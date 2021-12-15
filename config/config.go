package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Constants
}

type Constants struct {
	Host      string
	Port      int
	Databases []string
	Pager     string
}

func initConstants() (Constants, error) {
	viper.SetConfigName("config")
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		viper.AddConfigPath(fmt.Sprintf("%s/word", xdg))
	}
	viper.AddConfigPath("$HOME/.config/word")
	viper.AddConfigPath("$HOME/.word")
	err := viper.ReadInConfig()
	if err != nil {
		if _, notFound := err.(viper.ConfigFileNotFoundError); !notFound {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	viper.SetDefault("Host", "localhost")
	viper.SetDefault("Port", 2628)
	viper.SetDefault("Databases", []string{"gcide", "wn", "moby-thesaurus", "foldoc"})
	viper.SetDefault("Pager", "")

	var c Constants
	err = viper.Unmarshal(&c)
	return c, err
}

func New() (*Config, error) {
	config := Config{}
	constants, err := initConstants()
	config.Constants = constants
	if err != nil {
		return nil, err
	}
	return &config, err

}
