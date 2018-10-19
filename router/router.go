package router

import (
	"fmt"
	"github.com/leangeder/gravitywell/api"
	"github.com/leangeder/gravitywell/configuration"
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
		clusterConf := &configuration.ClusterConfig{
			Kind:       generalConfig.Kind,
			APIVersion: generalConfig.APIVersion,
		}
		api.ClusterApply(clusterConf)
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
