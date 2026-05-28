package validator

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const configName = ".commitlint-scope"

type Config struct {
	ScopeRegex *regexp.Regexp      `mapstructure:"scopeRegex"`
	Patterns   map[string][]string `mapstructure:"patterns"`
}

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(".")
	v.SetDefault("scopeRegex", regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); ok {
			return Config{}, nil
		}

		return Config{}, err
	}

	var cfg Config

	if err := v.Unmarshal(&cfg, regexDecode); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func regexDecode(cfg *mapstructure.DecoderConfig) {
	cfg.DecodeHook = func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.String && to == reflect.TypeOf(&regexp.Regexp{}) {
			val, ok := data.(string)
			if !ok {
				panic(fmt.Sprintf("expected string but got %T", data))
			}

			return regexp.Compile(val)
		}

		return data, nil
	}
}
