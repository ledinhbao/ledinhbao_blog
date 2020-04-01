package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Config is struct presentation for config.json file.
type Config map[string]interface{}

// ValueForKey takes key as string, dot-separated, return interface{}
func (c Config) ValueForKey(keyPath string) (interface{}, error) {
	keyArr := strings.Split(keyPath, ".")
	var value Config = c
	for i, key := range keyArr {
		if temp, ok := value[key]; ok {
			if i == len(keyArr)-1 {
				return temp, nil
			}
			if _, ok = temp.(map[string]interface{}); ok {
				value = temp.(map[string]interface{})
			} else {
				return nil, errors.New("Config-key doesn't exist")
			}
		} else {
			return nil, fmt.Errorf("Config key %s doesn't exist", key)
		}
	}
	return nil, errors.New("Undefined error")

}

// StringValueForKey takes key as string, dot-separated, return string,
// or error if value for intended key is not a string.
func (c Config) StringValueForKey(key string) (string, error) {
	value, err := c.ValueForKey(key)
	if err != nil {
		return "", err
	}
	if strValue, ok := value.(string); ok {
		return strValue, nil
	}
	return "", errors.New("Value for key is not a string")
}

// IntValueForKey takes key as string, dot-separated, return int,
// or error if value for intended key is not a int.
func (c Config) IntValueForKey(key string) (int, error) {
	value, err := c.ValueForKey(key)
	if err != nil {
		return 0, err
	}
	if strValue, ok := value.(int); ok {
		return strValue, nil
	}
	return 0, errors.New("Value for key is not a int")
}

// NewConfigFromJSONFile read from *.json and return Config object.
func NewConfigFromJSONFile(path string) (Config, error) {
	ifErr := func(err error) (Config, error) {
		return Config{}, err
	}

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return ifErr(err)
	}
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		return ifErr(err)
	}
	return config, nil
}
