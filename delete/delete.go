package delete

import (
	"context"
	"fmt"
	"log"

	"github.com/NissesSenap/pod-deleter/event"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodDeleter struct {
	kubeClient *kubernetes.Clientset
}

// SetupKubeClient
func SetupKubeClient(internal bool) (*PodDeleter, error) {

	// TODO add external case, there should be a interface
	if internal {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		// creates the clientset
		kubeClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return &PodDeleter{
			kubeClient: kubeClient,
		}, nil
	}
	return &PodDeleter{}, fmt.Errorf("No client")
}

// DeletePod, if not part of the criticalNamespaces the pod will be deleted
func (d *PodDeleter) DeletePod(falcoEvent event.Alert, criticalNamespaces map[string]bool) error {
	podName := falcoEvent.OutputFields.K8SPodName
	namespace := falcoEvent.OutputFields.K8SNsName
	log.Printf("PodName: %v & Namespace: %v", podName, namespace)

	log.Printf("Rule: %v", falcoEvent.Rule)
	if criticalNamespaces[namespace] {
		log.Printf("The pod %v won't be deleted due to it's part of the critical ns list: %v ", podName, namespace)
		return nil
	}

	log.Printf("Deleting pod %s from namespace %s", podName, namespace)
	err := d.kubeClient.CoreV1().Pods(namespace).Delete(context.Background(), podName, metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
