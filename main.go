package main

import (
	"log"
	"os"

	"github.com/NissesSenap/pod-deleter/delete"

	"github.com/NissesSenap/pod-deleter/event"
)

func main() {
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

	kubeClient, err := delete.SetupKubeClient(true)
	if err != nil {
		log.Fatalf("Unable to create in-cluster config: %v", err)
	}

	err = kubeClient.DeletePod(falcoEvent, criticalNamespaces)
	if err != nil {
		log.Fatalf("Unable to delete pod due to err %v", err)
	}
}
