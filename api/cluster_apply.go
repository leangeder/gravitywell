package api

import (
	"fmt"
	"github.com/leangeder/gravitywell/configuration"
)

func ClusterApply(config configuration.ClusterConfig) {
	fmt.Println("ClusterApply called")
	fmt.Println(config.Metadata.Name)
	fmt.Println(config.Spec.Replicas)
}
