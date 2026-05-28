package validator

import (
	"errors"
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

var ErrConfigRead = errors.New("error reading config")

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(".")
	v.SetDefault("scopeRegex", regexp.MustCompile(`^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`))

	var cfg Config

	if err := v.ReadInConfig(); err != nil {
		if _, ok := errors.AsType[viper.ConfigFileNotFoundError](err); !ok {
			return Config{}, fmt.Errorf("%w: %w", ErrConfigRead, err)
		}
	}

	if err := v.Unmarshal(&cfg, regexDecode); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

var ErrRegexDecode = errors.New("expected string for regexp decode")

func regexDecode(cfg *mapstructure.DecoderConfig) {
	cfg.DecodeHook = func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.String && to == reflect.TypeOf(&regexp.Regexp{}) {
			val, ok := data.(string)
			if !ok {
				return nil, fmt.Errorf("%w got %T", ErrRegexDecode, data)
			}

			return regexp.Compile(val)
		}

		return data, nil
	}
}
