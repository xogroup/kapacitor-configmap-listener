package configuration

import (
	"errors"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientSet Create a Kubernetes ClientSet based on flags
func NewClientSet(incluster *bool, kubeconfigpath *string) (*kubernetes.Clientset, error) {

	var config *rest.Config
	var err error

	if *incluster == true {
		config, err = NewInClusterConfig()
	} else {
		config, err = NewOutOfClusterConfig(kubeconfigpath)
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

// NewInClusterConfig Create an in-cluster configuration
func NewInClusterConfig() (*rest.Config, error) {

	return rest.InClusterConfig()
}

// NewOutOfClusterConfig Create an out-of-cluster configuration
func NewOutOfClusterConfig(kubeconfigpath *string) (*rest.Config, error) {

	if homepath := os.Getenv("HOME"); *kubeconfigpath == "" && homepath != "" {
		*kubeconfigpath = filepath.Join(homepath, ".kube", "config")
	} else {
		return nil, errors.New("absolute path required for kube config")
	}

	return clientcmd.BuildConfigFromFlags("", *kubeconfigpath)
}
