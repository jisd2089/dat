package settings

import (
	"io/ioutil"
	"os"

	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func CreateSettingsFromYAML(yamlPath string) (*DefaultSettings, error) {
	dataProvider, err := parseYAML(yamlPath)
	if err != nil {
		return nil, err
	}

	// return &DefaultSettings{yamlPath, dataProvider}, nil
	return &DefaultSettings{dataProvider}, nil
}

func parseYAML(yamlPath string) (map[interface{}]interface{}, error) {
	fd, err := os.Open(yamlPath)
	if err != nil {
		fmt.Printf("setting file not found:%s", yamlPath)
		return nil, fmt.Errorf("setting file not found:%s error:%s", yamlPath, err.Error())
	}

	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("reading setting:%s error:%s", yamlPath, err.Error())
	}

	receiver := make(map[interface{}]interface{})
	err = yaml.Unmarshal(content, &receiver)
	if err != nil {
		return nil, err
	}

	return receiver, nil
}
