package event

import (
	"testing"
)

var bodyData = "{\"output\":\"14:49:49.264147779: Notice A shell was spawned in a container with an attached terminal (user=root user_loginuid=-1 k8s.ns=default k8s.pod=alpine container=a15057582acc shell=sh parent=runc cmdline=sh -c uptime terminal=34816 container_id=a15057582acc image=alpine) k8s.ns=default k8s.pod=alpine container=a15057582acc k8s.ns=default k8s.pod=alpine container=a15057582acc\",\"priority\":\"Notice\",\"rule\":\"Terminal shell in container\",\"time\":\"2021-05-01T14:49:49.264147779Z\", \"output_fields\": {\"container.id\":\"a15057582acc\",\"container.image.repository\":\"alpine\",\"evt.time\":1619880589264147779,\"k8s.ns.name\":\"default\",\"k8s.pod.name\":\"alpine\",\"proc.cmdline\":\"sh -c uptime\",\"proc.name\":\"sh\",\"proc.pname\":\"runc\",\"proc.tty\":34816,\"user.loginuid\":-1,\"user.name\":\"root\"}}"

// TestRead
func TestReadWorking(t *testing.T) {
	bodyReqByte := []byte(bodyData)
	_, err := Read(bodyReqByte)
	if err != nil {
		t.Fatalf("The input probably needs to be cleaned up...")
	}
}

// TestReadWorkingBadData but correct syntax
func TestReadWorkingBadData(t *testing.T) {
	bodyReqByte := []byte("{\"data\": \"something\"}")
	_, err := Read(bodyReqByte)
	if err != nil {
		t.Fatalf("Got error %v", err)
	}
}

// TestReadBroken actually bad data
func TestReadBroken(t *testing.T) {
	bodyReqByte := []byte("something")
	_, err := Read(bodyReqByte)
	if err == nil {
		t.Fatalf("This should generate a error...")
	}
}

var criticalNamespaces = map[string]bool{
	"kube-system":     true,
	"kube-public":     true,
	"kube-node-lease": true,
	"falco":           true,
}

// TestCheckNamespaceTrue
func TestCheckNamespaceTrue(t *testing.T) {
	namespace := "not-critical"
	criticalNamespace := CheckNamespace(namespace, criticalNamespaces)
	if !criticalNamespace {
		t.Fatalf("Got namespace: %v, it shoulden't be critical", namespace)
	}
}

// TestCheckNamespaceFalse
func TestCheckNamespacFalse(t *testing.T) {
	namespace := "kube-system"
	criticalNamespace := CheckNamespace(namespace, criticalNamespaces)
	if criticalNamespace {
		t.Fatalf("Got namespace: %v, it should be critical", namespace)
	}
}
