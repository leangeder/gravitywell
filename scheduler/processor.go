package scheduler

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/platform"
	"github.com/leangeder/gravitywell/state"
	"github.com/leangeder/gravitywell/vcs"
	"k8s.io/api/core/v1"
)

func calculateRepositoryName(gitPath string) string {
	var extension = filepath.Ext(gitPath)
	var remoteVCSRepoName = gitPath[0 : len(gitPath)-len(extension)]
	splitStrings := strings.Split(remoteVCSRepoName, "/")
	return splitStrings[len(splitStrings)-1]
}

func ensureRepoExists(tempPath string, gitPath string, sshKeyPath string) (string, error) {
	repoName := calculateRepositoryName(gitPath)
	if _, err := os.Stat(path.Join(tempPath, repoName)); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("Fetching repository %s into %s\n", repoName, path.Join(tempPath, repoName)))
		gvcs := new(vcs.GitVCS)
		_, err = vcs.Fetch(gvcs, path.Join(tempPath, repoName), gitPath, sshKeyPath)
		if err != nil {
			return repoName, err
		}
	} else {
		log.Debug(fmt.Sprintf("Using existing repository %s", path.Join(tempPath, repoName)))
	}
	return repoName, nil
}

func process(opt configuration.Options, cluster configuration.Application) *state.Capture {

	stateCapture := &state.Capture{
		ClusterName:     cluster.Name,
		DeploymentState: make(map[string]state.Details),
	}

	log.Info(fmt.Sprintf("Switching to cluster: %s\n", cluster.Name))
	restclient, k8siface, err := platform.GetKubeClient(cluster.Name)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	for _, deployment := range cluster.Deployments {
		log.Info(fmt.Sprintf("Loading deployment %s\n", deployment.Deployment.Name))

		remoteVCSRepoName, err := ensureRepoExists(opt.TempVCSPath, deployment.Deployment.Git, opt.SSHKeyPath)
		if err != nil {
			log.Error(err.Error())
			stateCapture.DeploymentState[deployment.Deployment.Name] = state.Details{State: state.EDeploymentStateError}
			return stateCapture
		}
		//---------------------------------
		for _, a := range deployment.Deployment.Action {
			if a.Execute.Shell != "" {
				log.Info(fmt.Sprintf("Running shell command %s\n", a.Execute.Shell))
				if err := ShellCommand(a.Execute.Shell, path.Join(opt.TempVCSPath, remoteVCSRepoName), true); err != nil {
					log.Error(err.Error())
				}
			}
			//---------------------------------
			var commandFlag configuration.CommandFlag
			if a.Execute.Kubectl.Command == "" {
				log.Debug("No Kubernetes action to run aborting (supports: create/apply/replace)")
				continue
			}
			switch strings.ToLower(a.Execute.Kubectl.Command) {
			case "apply":
				log.Println("Using apply command")
				commandFlag = configuration.Apply
			case "create":
				log.Println("Using create command")
				commandFlag = configuration.Create
			case "replace":
				log.Println("Using replace command")
				commandFlag = configuration.Replace
			default:

			}
			//---------------------------------
			fileList := []string{}
			err := filepath.Walk(path.Join(opt.TempVCSPath, remoteVCSRepoName, a.Execute.Kubectl.Path), func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
			if err != nil {
				log.Error(err.Error())

			}
			for _, file := range fileList {
				log.Debug(fmt.Sprintf("Attempting to deploy %s\n", file))
				if _, err = os.Stat(file); os.IsNotExist(err) {
					continue
				}
				if sa, _ := os.Stat(file); sa.IsDir() {
					continue
				}
				var stateResponse state.State
				//---------------------------------
				log.Debug(fmt.Sprintf("Running..."))
				if deployment.Deployment.CreateNamespace {
					log.Debug("Creating namespace...")
					nsclient := k8siface.CoreV1().Namespaces()
					cm := &v1.Namespace{}
					cm.Name = deployment.Deployment.Namespace
					cm.Namespace = deployment.Deployment.Namespace
					_, err := nsclient.Create(cm)
					if err != nil {
						log.Error(fmt.Sprintf("Could not deploy Namespace resource %s due to %s", cm.Name, err.Error()))
					}
				}
				if stateResponse, err = platform.DeployFromFile(restclient, k8siface, file, deployment.Deployment.Namespace, opt, commandFlag); err != nil {
					log.Error(err.Error())
				}
				var output = ""
				var hasError = false
				if err != nil {
					output = fmt.Sprintf("File: %s Namespace :%s Error: %s", file, deployment.Deployment.Namespace, err)
					hasError = true
				} else {
					output = fmt.Sprintf("File: %s Namespace :%s", file, deployment.Deployment.Namespace)
				}
				stateCapture.DeploymentState[deployment.Deployment.Name] = state.Details{State: stateResponse, HasDetail: true,
					Detail: output, HasError: hasError}
			}
			//---------------------------------
		}
	}
	return stateCapture
}
