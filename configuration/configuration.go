package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type GeneralConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Metadata   interface{} `yaml:"Metadata"`
	Spec       interface{} `yaml:"Spec"`
}

func NewConfiguration(conf string) (*GeneralConfig, error) {
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
