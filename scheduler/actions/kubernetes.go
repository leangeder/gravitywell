package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform"
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func clientForCluster(clusterName string) (*rest.Config, kubernetes.Interface) {
	log.Info(fmt.Sprintf("Switching to cluster: %s\n", clusterName))
	restclient, k8siface, err := platform.GetKubeClient(clusterName)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	return restclient, k8siface
}

func ExecuteKubernetesAction(action kinds.Action, clusterName string,
	deployment kinds.Application,
	commandFlag configuration.CommandFlag, opt configuration.Options, repoName string) {
	var deploymentPath = "."
	shouldAwaitDeployment := false
	if tp, ok := action.Execute.Configuration["Path"]; ok && tp != "" {
		deploymentPath = tp
	}
	if tp, ok := action.Execute.Configuration["AwaitDeployment"]; ok && tp != "" {
		b, err := strconv.ParseBool(tp)
		if err != nil {
			log.Error(err.Error())
		}
		shouldAwaitDeployment = b
	}

	fileList := []string{}
	err := filepath.Walk(path.Join(opt.TempVCSPath,
		repoName, deploymentPath),
		func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return nil
		})
	if err != nil {
		log.Error(err.Error())
	}
	_, k8siface := clientForCluster(clusterName)
	err = platform.GenerateDeploymentPlan(
		k8siface, fileList,
		deployment.Namespace, opt, commandFlag, shouldAwaitDeployment)
	if err != nil {
		log.Error(err.Error())
	}
	//---------------------------------
}
