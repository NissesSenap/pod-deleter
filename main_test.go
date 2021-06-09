// +build integration

// when using build integration package main dosen't seem to get imported correctly...
package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func k8sClient() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

const bodyData string = "{\"output\":\"14:49:49.264147779: Notice A shell was spawned in a container with an attached terminal (user=root user_loginuid=-1 k8s.ns=default k8s.pod=alpine container=a15057582acc shell=sh parent=runc cmdline=sh -c uptime terminal=34816 container_id=a15057582acc image=alpine) k8s.ns=default k8s.pod=alpine container=a15057582acc k8s.ns=default k8s.pod=alpine container=a15057582acc\",\"priority\":\"Notice\",\"rule\":\"Terminal shell in container\",\"time\":\"2021-05-01T14:49:49.264147779Z\", \"output_fields\": {\"container.id\":\"a15057582acc\",\"container.image.repository\":\"alpine\",\"evt.time\":1619880589264147779,\"k8s.ns.name\":\"default\",\"k8s.pod.name\":\"alpine\",\"proc.cmdline\":\"sh -c uptime\",\"proc.name\":\"sh\",\"proc.pname\":\"runc\",\"proc.tty\":34816,\"user.loginuid\":-1,\"user.name\":\"root\"}}"

func TestIntegrationAllowNS(t *testing.T) {
	ctx := context.Background()
	var (
		namespace string = "default"
		podName   string = "alpine"
		image     string = "alpine"
	)
	cmd := []string{"sh", "-c", "sleep 600"}
	kubeClient, err := k8sClient()
	if err != nil {
		t.Fatalf("Unable to create a kubernetes client: %v", err)
	}

	err = os.Setenv(envBlockNamespace, "super-critical")
	if err != nil {
		t.Fatalf("Unable to set env : %v", err)
	}

	err = os.Setenv(envBodyReq, bodyData)
	if err != nil {
		t.Fatalf("Unable to set env : %v", err)
	}

	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      podName, // TODO start using generate output to be able to run test at the same time.
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    podName,
					Image:   image,
					Command: cmd,
				},
			},
		},
	}

	// use the output to get the pod name from the generation
	_, err = kubeClient.CoreV1().Pods(namespace).Create(ctx, pod, metaV1.CreateOptions{})
	if err != nil {
		t.Fatalf("Unable to create pod %v in namespace %v err: %v", podName, namespace, err)
	}

	// There should a wait command... instead of this sleep
	// probably https://pkg.go.dev/k8s.io/apimachinery/pkg/util/wait#PollImmediateUntil but I was hoping for something easier
	time.Sleep(2 * time.Second)
	// run main
	main()
	dyingPod, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err == nil {
		if dyingPod.GetObjectMeta().GetDeletionTimestamp() == nil {
			t.Fatalf("The pod %v should be dying or allready be gone in namespace: %v", podName, namespace)
		}
	}
}
