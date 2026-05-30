package validator

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const ConfigName = ".commitlint-scope.yaml"

type PatternItem struct {
	Scopes []string `koanf:"scopes"`
	Files  []string `koanf:"files"`
}

type Config struct {
	ScopeRegex *regexp.Regexp `koanf:"scopeRegex"`
	Patterns   []PatternItem  `koanf:"patterns"`
}

var (
	ErrConfigRead  = errors.New("error reading config")
	ErrRegexDecode = errors.New("expected string for regexp decode")
)

func LoadConfig() (Config, error) {
	k := koanf.New(".")

	defaultRegex := `^[a-z]+(?:\((?P<scope>[^)]+)\))?!?:\s`
	_ = k.Set("scopeRegex", defaultRegex)

	if err := k.Load(file.Provider(ConfigName), yaml.Parser()); err != nil {
		if !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("%w: %w", ErrConfigRead, err)
		}
	}

	var cfg Config

	unmarshalConf := koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			Result:     &cfg,
			DecodeHook: regexDecodeHook,
		},
	}

	if err := k.UnmarshalWithConf("", &cfg, unmarshalConf); err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrConfigRead, err)
	}

	return cfg, nil
}

func regexDecodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf(&regexp.Regexp{}) {
		val, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("%w got %T", ErrRegexDecode, data)
		}

		return regexp.Compile(val)
	}

	return data, nil
}
