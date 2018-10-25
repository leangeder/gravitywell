package platform

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/state"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Beta1DaemonSetResouce(k kubernetes.Interface, sts *v1beta1.DaemonSet, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found statefulset resource")
	dsclient := k.Extensions().DaemonSets(namespace)

	if opts.DryRun {
		_, err := dsclient.Get(sts.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", sts.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", sts.Name))
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
			log.Error(fmt.Sprintf("Could not deploy sts resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := dsclient.Create(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy sts resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := dsclient.UpdateStatus(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not apply Statefulset resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Statefulset updated")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")

}
