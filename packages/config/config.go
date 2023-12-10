package config

import (
	"fmt"
	"errors"
	"encoding/json"
	"io/ioutil"
    "gopkg.in/yaml.v3"
)

type GlobalSettings struct {
	LoggingLevel int `yaml:"loggingLevel" json:"loggingLevel"`
	Port int `yaml:"port" json:"port"`
}

var (
	Settings GlobalSettings
)

func ParseConfig(config string) error {
	var errs [4]error

	Settings, errs[0] = parseYamlConfig(config)
	if errs[0] == nil {
		// if the Port value is still 0, we either didn't
		// a good parse, or we didn't get any errors. Either
		// way this is an issue
		if Settings.Port == 0 {
			errs[1] = fmt.Errorf("%s", "After unmarshalling yaml, port is still 0 [unknown error]")
		} else {
			return nil
		}
	}

	Settings, errs[2] = parseJsonConfig(config)
	if errs[2] == nil {
		// just like in the yaml parsing. this is a misconfigured
		// file - no matter what format it is in
		if Settings.Port == 0 {
			errs[3] = fmt.Errorf("%s", "After unmarshalling json, port is still 0 [unknown error]")
		} else {
			return nil
		}
	}

	return fmt.Errorf("%s", errors.Join(errs[:]...))
}

type unmarshalFunc func([]byte, any) error

func parseConfigFile(file string, unmarshal unmarshalFunc) (GlobalSettings, error) {
	var conf GlobalSettings

    fileContents, err := ioutil.ReadFile(file)
    if err != nil {
        return conf, err
    }

    err = unmarshal(fileContents, &conf)
    if err != nil {
        return conf, err
    }

    return conf, nil
}
 
func parseYamlConfig(file string) (GlobalSettings, error) {
    return parseConfigFile(file, yaml.Unmarshal)
}

func parseJsonConfig(file string) (GlobalSettings, error) {
    return parseConfigFile(file, json.Unmarshal)
}
