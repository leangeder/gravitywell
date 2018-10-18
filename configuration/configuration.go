package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Cluster struct {
	Replicas int `yaml:"replicas"`
}

type ClusterDeployment struct {
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
type Configuration struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Spec       interface{} `yaml:"Spec"`
}

type ClusterConfiguration struct {
	Kind       string  `yaml:"Kind"`
	APIVersion string  `yaml:"APIVersion"`
	Spec       Cluster `yaml:"Spec"`
}

type ClusterDeploymentConfiguration struct {
	Kind       string            `yaml:"Kind"`
	APIVersion string            `yaml:"APIVersion"`
	Spec       ClusterDeployment `yaml:"Spec"`
}

//NewConfiguration creates a deserialised yaml object
func NewConfiguration(conf string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	c := Configuration{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("Failed to validate syntax: %s", conf)
		return nil, err
	}

	if c.Kind == "Cluster" {
		c.Spec = c.Spec.(Cluster)
	} else if c.Kind == "ClusterDeployment" {
		c.Spec = c.Spec.(ClusterDeployment)
	}
	return &c, nil
}
