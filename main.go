package main

import (
	"context"
	"log"
	"os"

	"github.com/NissesSenap/pod-deleter/event"
	"github.com/NissesSenap/pod-deleter/poddelete"
)

func main() {
	ctx := context.Background()
	blockNamespaces := map[string]bool{
		"kube-system":     true,
		"kube-public":     true,
		"kube-node-lease": true,
		"falco":           true,
	}

	var allowNamespaces map[string]bool
	// TODO can i check if alllowNamepsaces have done make and from that draw the same conclusion?
	var allowList bool

	blockNamespaceEnv := os.Getenv("BLOCK_NAMESPACE")
	if blockNamespaceEnv != "" {
		blockNamespaces = event.AddItemsToHashMap(blockNamespaceEnv, blockNamespaces)
	}

	allowNamespaceEnv := os.Getenv("ALLOW_NAMESPACE")

	if blockNamespaceEnv != "" && allowNamespaceEnv != "" {
		log.Fatalf("Both env BLOCK_NAMESPACE: %v & ALLOW_NAMESPACE: %v, can't be defined", blockNamespaceEnv, allowNamespaceEnv)
	}

	if allowNamespaceEnv != "" {
		allowNamespaces = make(map[string]bool)
		allowList = true
		allowNamespaces = event.AddItemsToHashMap(allowNamespaceEnv, allowNamespaces)
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

	var checkNamespace bool
	if allowList {
		checkNamespace = event.CheckAllowNamespace(namespace, allowNamespaces)
	} else {
		checkNamespace = event.CheckBlockNamespace(namespace, blockNamespaces)
	}

	if !checkNamespace {
		// TODO, we should add what namespace that blockNamespace output for easy debug
		log.Printf("Not deleting pod: %v, in namespace: %v", podName, namespace)
		os.Exit(0)
	}

	kubeClient, err := poddelete.SetupKubeClient(false)
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
