package platform

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/leangeder/gravitywell/configuration"
	"github.com/leangeder/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	v1betav1 "k8s.io/api/extensions/v1beta1"
	v1polbeta "k8s.io/api/policy/v1beta1"
	v1rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//This is required for gcp auth provider scope
)

// GetKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func GetKubeClient(context string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string) (*rest.Config, error) {
	config, err := getConfig(context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getConfig returns a Kubernetes client config for a given context.
func getConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

//DeployFromFile ...
func DeployFromFile(config *rest.Config, k kubernetes.Interface, path string, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	f, err := os.Open(path)
	if err != nil {
		return state.EDeploymentStateError, err
	}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return state.EDeploymentStateError, err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, _ := decode(raw, nil, nil)

	log.Debug(fmt.Sprintf("%++v\n\n", obj.GetObjectKind()))

	var response state.State
	var e error
	switch obj.(type) {
	case *v1beta1.Deployment:
		response, e = execV1Beta1DeploymentResouce(k, obj.(*v1beta1.Deployment), namespace, opts, commandFlag)
	case *v1beta2.Deployment:
		response, e = execV1Beta2DeploymentResouce(k, obj.(*v1beta2.Deployment), namespace, opts, commandFlag)
	case *v1beta1.StatefulSet:
		response, e = execV1Beta1StatefulSetResouce(k, obj.(*v1beta1.StatefulSet), namespace, opts, commandFlag)
	case *v1.Service:
		response, e = execV1ServiceResouce(k, obj.(*v1.Service), namespace, opts, commandFlag)
	case *v1.ConfigMap:
		response, e = execV1ConfigMapResouce(k, obj.(*v1.ConfigMap), namespace, opts, commandFlag)
	case *v1polbeta.PodDisruptionBudget:
		response, e = execV1Beta1PodDisruptionBudgetResouce(k, obj.(*v1polbeta.PodDisruptionBudget), namespace, opts, commandFlag)
	case *v1.ServiceAccount:
		response, e = execV1ServiceAccountResouce(k, obj.(*v1.ServiceAccount), namespace, opts, commandFlag)
	//V1 RBAC
	case *v1rbac.ClusterRoleBinding:
		response, e = execV1RbacClusterRoleBindingResouce(k, obj.(*v1rbac.ClusterRoleBinding), namespace, opts, commandFlag)
	case *v1rbac.Role:
		response, e = execV1RbacRoleResouce(k, obj.(*v1rbac.Role), namespace, opts, commandFlag)
	case *v1rbac.RoleBinding:
		response, e = exec1VRbacRoleBindingResouce(k, obj.(*v1rbac.RoleBinding), namespace, opts, commandFlag)
	case *v1rbac.ClusterRole:
		response, e = execV1AuthClusterRoleResouce(k, obj.(*v1rbac.ClusterRole), namespace, opts, commandFlag)
	case *v1betav1.DaemonSet:
		response, e = execV1Beta1DaemonSetResouce(k, obj.(*v1betav1.DaemonSet), namespace, opts, commandFlag)
	default:
		log.Error("Unable to convert API resource")
	}

	return response, e
}
