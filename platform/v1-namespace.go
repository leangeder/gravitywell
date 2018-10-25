package platform

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/state"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1NamespaceResouce(k kubernetes.Interface, ss *v1.Namespace, _ string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found namespace resource")
	ssclient := k.CoreV1().Namespaces()

	if opts.DryRun {
		_, err := ssclient.Get(ss.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Namespace resource %s does not exist\n", ss.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Namespace resource %s exists\n", ss.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		ssclient.Delete(ss.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := ssclient.Create(ss)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Namespace resource %s due to %s", ss.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Namespace deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := ssclient.Create(ss)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Namespace resource %s due to %s", ss.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Namespace deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := ssclient.Update(ss)
		if err != nil {
			log.Error(fmt.Sprintf("Could not apply Namespace resource %s due to %s", ss.Name, err.Error()))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Namespace updated")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
