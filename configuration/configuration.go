package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Application struct {
	Name        string `yaml:"Name"`
	Deployments []struct {
		Deployment struct {
			Name            string `yaml:"Name"`
			Namespace       string `yaml:"Namespace"`
			CreateNamespace bool   `yaml:"CreateNamespace"`
			Git             string `yaml:"Git"`
			Action          []struct {
				Execute struct {
					Shell   string `yaml:"Shell"`
					Kubectl struct {
						Path    string `yaml:"Path"`
						Type    string `yaml:"Type"`
						Command string `yaml:"Command"`
					} `yaml:"Kubectl"`
				} `yaml:"Execute"`
			} `yaml:"Action"`
		} `yaml:"Deployment"`
	} `yaml:"Deployments"`
}

//Configuration ...
type GeneralConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Spec       interface{} `yaml:"Spec"`
}

type ApplicationConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Spec       Application `yaml:"Spec"`
}

//NewConfiguration creates a deserialised yaml object
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
