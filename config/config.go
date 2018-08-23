package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config of autofresh
type Config struct {
	Version    bool
	Hidebanner bool
	Watchman   string
	Build      string
	Run        string
	Suffixes   []string
}

// LoadConfig takes in a cobra Command and unmarshalls it into the custom config
// struct. It also binds pflags to viper, sets environment variables, and reads
// in a autofresh-config.{json,yaml,toml}.
func LoadConfig(cmd *cobra.Command) (Config, []error) {
	var errs []error
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		errs = append(errs, fmt.Errorf("failed to bind flags, error: %s", err.Error()))
	}

	viper.SetEnvPrefix("AUTOFRESH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigFile("autofresh-config.json")
	if err := viper.ReadInConfig(); err != nil {
		errs = append(errs, fmt.Errorf("failed to read config file, error: %s", err.Error()))
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Unable to unmarshal into struct, error: %s\n", err.Error())
	}
	return conf, errs
}
