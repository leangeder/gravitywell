package platform

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/state"
	"k8s.io/api/batch/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1JobSetResouce(k kubernetes.Interface, sts *v1.Job, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found Job resource")
	dsclient := k.Batch().Jobs(namespace)

	if opts.DryRun {
		_, err := dsclient.Get(sts.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Job resource %s does not exist\n", sts.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Job resource %s exists\n", sts.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		dsclient.Delete(sts.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := dsclient.Create(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Job resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Job deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := dsclient.Create(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Job resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Job deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := dsclient.UpdateStatus(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not apply Job resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Job updated")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")

}
