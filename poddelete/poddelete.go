package poddelete

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PodDeleter struct {
	kubeClient kubernetes.Interface
}

const falcoAnnotation = "falco.org/protected"

// SetupKubeClient
func SetupKubeClient(internal bool) (*PodDeleter, error) {

	var config *rest.Config
	var err error

	// Check if we are inside the Kubernetes cluster using a ServiceAccount
	if internal {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		// TODO Should probably migrate to a input variable instead?
		kubeconfig := os.Getenv("KUBECONFIG")

		if kubeconfig == "" {
			kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return &PodDeleter{}, err
		}
	}
	// creates the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &PodDeleter{kubeClient: kubeClient}, nil
}

// CheckPodAnnotation looks if pod is protected by annotation falco.org/protected: "true"
func (d *PodDeleter) CheckPodAnnotation(ctx context.Context, namespace, podName string) (bool, error) {

	podData, err := d.kubeClient.CoreV1().Pods(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("Unable to get pod: %v in namespace: %v, due to error: %v", podName, namespace, err)
	}

	if podData.Annotations[falcoAnnotation] == "true" || podData.Annotations[falcoAnnotation] == "True" {
		return true, nil
	}
	return false, nil
}

// DeletePod, if not part of the criticalNamespaces the pod will be deleted
func (d *PodDeleter) DeletePod(ctx context.Context, namespace, podName string) error {

	err := d.kubeClient.CoreV1().Pods(namespace).Delete(ctx, podName, metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
