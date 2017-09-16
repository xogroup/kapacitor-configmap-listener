package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubeClient create a Kubernetes ClientSet based on flags
func NewKubeClient(inCluster *bool, kubeConfigPath *string) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error

	if *inCluster == true {
		config, err = NewInClusterConfig()
	} else {
		config, err = NewOutOfClusterConfig(kubeConfigPath)
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// NewInClusterConfig creates an in-cluster configuration
func NewInClusterConfig() (*rest.Config, error) {

	return rest.InClusterConfig()
}

// NewOutOfClusterConfig create an out-of-cluster configuration.  If the `kubeconfigpath` is empty
// an attempt is made to locate the `HOME` directory and a path for `~/.kube/config` is used as default.
func NewOutOfClusterConfig(kubeConfigPath *string) (*rest.Config, error) {

	if homePath := os.Getenv("HOME"); *kubeConfigPath == "" && homePath != "" {
		*kubeConfigPath = filepath.Join(homePath, ".kube", "config")
	} else {
		return nil, fmt.Errorf("absolute path required for kube config")
	}

	return clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
}
