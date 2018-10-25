package router

import (
	"errors"
	"fmt"
	"github.com/leangeder/gravitywell/api"
	"github.com/leangeder/gravitywell/configuration"
	yaml "gopkg.in/yaml.v2"
)

type RouterPath struct {
	verb string
	kind string
}

func Run(verb string, generalConfig *configuration.GeneralConfig, opt *configuration.Options) {

	routerPath := &RouterPath{verb: verb, kind: generalConfig.Kind}
	fmt.Println(verb, generalConfig.Kind)

	switch *routerPath {
	case RouterPath{verb: "apply", kind: "Cluster"}:
		fmt.Println("Apply cluster")
		conf, err := MapToClusterConfig(generalConfig)
		if err != nil {
			fmt.Println("Error mapping clusterConfig:", err.Error())
			return
		}
		api.ClusterApply(*conf)
	case RouterPath{verb: "apply", kind: "Application"}:
		fmt.Println("Apply Application")
		conf, err := MapToApplicationConfig(generalConfig)
		if err != nil {
			fmt.Println("Error mapping applicationConfig:", err.Error())
			return
		}
		api.ApplicationApply(conf, opt)
	default:
		fmt.Println("Route not recognize")
	}
}

func MapToApplicationConfig(conf *configuration.GeneralConfig) (*configuration.ApplicationConfig, error) {
	bytes, err := yaml.Marshal(conf.Strategy)
	if err != nil {
		return nil, errors.New("Failed to marshal Strategy")
	}
	strategy := configuration.Strategy{}
	err = yaml.Unmarshal(bytes, &strategy)
	if err != nil {
		return nil, errors.New("Failed to unmarshal Strategy")
	}

	appConf := &configuration.ApplicationConfig{
		Kind:       conf.Kind,
		APIVersion: conf.APIVersion,
		Strategy:   strategy,
	}
	return appConf, nil
}

func MapToClusterConfig(conf *configuration.GeneralConfig) (*configuration.ClusterConfig, error) {
	bytes, err := yaml.Marshal(conf.Spec)
	if err != nil {
		return nil, errors.New("Failed to marshal Spec")
	}
	spec := configuration.ClusterSpec{}
	err = yaml.Unmarshal(bytes, &spec)
	if err != nil {
		return nil, errors.New("Failed to unmarshal Spec")
	}

	bytes, err = yaml.Marshal(conf.Metadata)
	if err != nil {
		return nil, errors.New("Failed to marshal Metadata")
	}
	metadata := configuration.ObjectMeta{}
	err = yaml.Unmarshal(bytes, &metadata)
	if err != nil {
		return nil, errors.New("Failed to unmarshal Metadata")
	}

	clusterConf := &configuration.ClusterConfig{
		Kind:       conf.Kind,
		Metadata:   metadata,
		APIVersion: conf.APIVersion,
		Spec:       spec,
	}
	return clusterConf, nil
}