package main

import (
	"context"
	"log"
	"os"

	"github.com/NissesSenap/pod-deleter/delete"

	"github.com/NissesSenap/pod-deleter/event"
)

func main() {
	ctx := context.Background()
	criticalNamespaces := map[string]bool{
		"kube-system":     true,
		"kube-public":     true,
		"kube-node-lease": true,
		"falco":           true,
	}

	bodyReq := os.Getenv("BODY")
	if bodyReq == "" {
		log.Fatalf("Need to get environment variable BODY")
	}
	bodyReqByte := []byte(bodyReq)
	falcoEvent, err := event.Read(bodyReqByte)
	if err != nil {
		log.Fatalf("The data doesent match the struct %v", err)
	}

	podName := falcoEvent.OutputFields.K8SPodName
	namespace := falcoEvent.OutputFields.K8SNsName

	criticalNamespace := event.CheckNamespace(namespace, criticalNamespaces)

	if !criticalNamespace {
		// TODO, we should add what namespace that criticalNamespace output for easy debug
		log.Printf("Not deleting pod: %v, in namespace: %v due to it's in criticalNamepsace", namespace, podName)
		os.Exit(0)
	}

	kubeClient, err := delete.SetupKubeClient(false)
	if err != nil {
		log.Fatalf("Unable to create in-cluster config: %v", err)
	}

	podProtected, err := kubeClient.CheckPodAnnotation(ctx, namespace, podName)
	if err != nil {
		log.Fatal("Unable to get pod annotation: ", err)
	}
	if podProtected {
		log.Printf("Not deleting pod: %v, in namespace: %v due to it got the annotation falco.org/protected: %v", namespace, podName, podProtected)
		os.Exit(0)
	}

	err = kubeClient.DeletePod(ctx, namespace, podName)
	if err != nil {
		log.Fatalf("Unable to delete pod due to err %v", err)
	}
	log.Printf("Deleted pod: %v in namespace: %v", podName, namespace)
}
