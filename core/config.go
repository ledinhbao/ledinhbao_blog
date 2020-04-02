package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// ConfigIndexError presents key does not exist.
type ConfigIndexError struct {
	Key string
}

// Config is struct presentation for config.json file.
type Config map[string]interface{}

func (err ConfigIndexError) Error() string {
	return fmt.Sprintf("Cound not find values for keys %s", err.Key)
}

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
				return nil, ConfigIndexError{Key: key}
			}
		} else {
			return nil, ConfigIndexError{Key: key}
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
	if res, ok := value.(int); ok {
		return res, nil
	}
	return 0, errors.New("Value for key is not a int")
}

// ConfigValueForKey takes key as string, dot-separated, return Config object,
// or error if value for intended key cannot form a Config object.
func (c Config) ConfigValueForKey(key string) (Config, error) {
	value, err := c.ValueForKey(key)
	if err != nil {
		return Config{}, err
	}
	if _, ok := value.(map[string]interface{}); ok {
		var res Config = value.(map[string]interface{})
		return res, nil
	}
	return Config{}, errors.New("Value for key " + key + "is single field value")
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
