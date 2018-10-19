package router

import (
	"fmt"
	"github.com/leangeder/gravitywell/api"
	"github.com/leangeder/gravitywell/configuration"
	yaml "gopkg.in/yaml.v2"
	"log"
)

type RouterPath struct {
	verb string
	kind string
}

func Run(verb string, generalConfig *configuration.GeneralConfig) {

	routerPath := &RouterPath{verb: verb, kind: generalConfig.Kind}
	fmt.Println(verb, generalConfig.Kind)

	switch *routerPath {
	case RouterPath{verb: "apply", kind: "Cluster"}:
		fmt.Println("Apply cluster")
		bytes, err := yaml.Marshal(generalConfig.Spec)
		if err != nil {
			log.Printf("Failed to marshal Spec")
			return
		}
		spec := configuration.ClusterSpec{}
		err = yaml.Unmarshal(bytes, &spec)
		if err != nil {
			log.Printf("Failed to unmarshal Spec")
			return
		}

		clusterConf := &configuration.ClusterConfig{
			Kind:       generalConfig.Kind,
			APIVersion: generalConfig.APIVersion,
			Spec:       spec,
		}
		api.ClusterApply(*clusterConf)
	case RouterPath{verb: "apply", kind: "Application"}:
		fmt.Println("Apply Application")
		appConf := &configuration.ApplicationConfig{
			Kind:       generalConfig.Kind,
			APIVersion: generalConfig.APIVersion,
		}
		api.ApplicationApply(appConf)
	default:
		fmt.Println("Route not recognize")
	}
}
