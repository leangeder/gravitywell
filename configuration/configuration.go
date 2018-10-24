package configuration

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type GeneralConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Metadata   interface{} `yaml:"Metadata"`
	Spec       interface{} `yaml:"Spec"`
	Strategy   interface{} `yaml:"Strategy"`
}

func NewConfiguration(conf string) ([]*GeneralConfig, error) {
	path, err := os.Stat(conf)

	if err != nil {
		log.Printf("Failed to identify status of conf")
		return nil, err
	}

	switch mode := path.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(conf)
		if err != nil {
			return nil, err
		}

		configFiles := make([]*GeneralConfig, len(files))

		for i, item := range files {
			generalConf, err := GenerateFileToConf(conf + item.Name())
			if err != nil {
				return nil, errors.New("Invalid file: " + item.Name())
			}
			configFiles[i] = generalConf
		}
		return configFiles, err
	case mode.IsRegular():
		generalConf, err := GenerateFileToConf(conf)
		return []*GeneralConfig{generalConf}, err
	}
	return nil, errors.New("Config is not a file or a directory")
}

func GenerateFileToConf(conf string) (*GeneralConfig, error) {
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	c := GeneralConfig{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("Failed to validate syntax: %s", conf)
		return nil, err
	}

	return &c, nil
}
