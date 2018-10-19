package api

import (
	"fmt"
	"github.com/leangeder/gravitywell/configuration"
)

func ClusterApply(config configuration.ClusterConfig) {
	fmt.Println("ClusterApply called")
	fmt.Println(config.Kind)
	fmt.Println(config.Metadata.Name)
	fmt.Println(config.Spec.Replicas)
	fmt.Println(config.Spec.Template.Spec.ServiceAccountName)
	fmt.Println(config.Spec.Template.Spec.Providers)

	fmt.Println(config.Spec.Template.Spec.Providers[0])

	fmt.Println(config.Spec.Template.Spec.Providers[0].Name)
	fmt.Println(config.Spec.Template.Spec.Providers[0].NodePools)
	fmt.Println(config.Spec.Template.Spec.Providers[0].NodePools[0].Name)
}
