package wikipedia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Category    CategoryFilter
	TitleLength TitleLengthFilter
	Title       TitleFilter
}

func ReadConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("Can't open config file `%s` `%s`", path, err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Can't parse config `%s` `%s`", path, err)
	}

	return config, nil
}
