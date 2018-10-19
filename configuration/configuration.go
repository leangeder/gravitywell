package configuration

import (
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
}

func NewConfiguration(conf string) ([]*GeneralConfig, error) {

	file, err := os.Stat(conf)

	if err != nil {
		log.Printf("Failed to identify status of conf")
		return nil, err
	}

	switch mode := file.Mode(); {
	case mode.IsDir():
		listFile, err := filepath.Glob(conf + "*")
		if err != nil {
			return nil, err
		}

		for _, item := range listFile {
			generalConf, err := GenerateFileToConf(conf + item)
		}
		return []*GeneralConfig{generalConf}, err
	case mode.IsRegular():
		generalConf, err := GenerateFileToConf(conf)
		return []*GeneralConfig{generalConf}, err
	}
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
