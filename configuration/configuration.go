package configuration

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
		log.Error("Failed to identify status of conf")
		return nil, err
	}

	switch mode := path.Mode(); {
	case mode.IsDir():
		files, err := ioutil.ReadDir(conf)
		if err != nil {
			return nil, err
		}

		var configFiles []*GeneralConfig

		for _, item := range files {
			generalConf, err := GenerateFileToConf(conf + item.Name())
			if err != nil {
				return nil, errors.New("Invalid file: " + item.Name())
			}
			configFiles = append(configFiles, generalConf...)
		}
		return configFiles, err
	case mode.IsRegular():
		return GenerateFileToConf(conf)
	}
	return nil, errors.New("Config is not a file or a directory")
}

func GenerateFileToConf(conf string) ([]*GeneralConfig, error) {
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	/*err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("Failed to validate syntax: %s", conf)
		return nil, err
	}
	*/
	var configFiles []*GeneralConfig
	dec := yaml.NewDecoder(strings.NewReader(string(bytes)))
	for {
		log.Debug("Document on file " + conf)
		value := GeneralConfig{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		configFiles = append(configFiles, &value)
	}
	return configFiles, nil
}
