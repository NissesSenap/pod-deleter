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

// TestCheckAllowNamespace
func TestCheckAllowNamespace(t *testing.T) {
	data := []struct {
		namespace      string
		namespaces     map[string]bool
		expectedOutput bool
	}{
		{
			namespace: "in-allow-list",
			namespaces: map[string]bool{
				"ns1":           true,
				"ns2":           true,
				"in-allow-list": true,
			},
			expectedOutput: true,
		},
		{
			namespace: "not-in-allow-list",
			namespaces: map[string]bool{
				"ns1": true,
				"ns2": true,
			},
			expectedOutput: false,
		},
	}
	for _, single := range data {
		t.Run("", func(single struct {
			namespace      string
			namespaces     map[string]bool
			expectedOutput bool
		}) func(t *testing.T) {
			return func(t *testing.T) {
				output := CheckAllowNamespace(single.namespace, single.namespaces)
				if output != single.expectedOutput {
					t.Errorf("Got: %v, expected %v, for namespace %v", output, single.expectedOutput, single.namespace)
				}
			}
		}(single))
	}
}

// TestCheckBlockNamespace
func TestCheckBlockNamespace(t *testing.T) {
	data := []struct {
		namespace      string
		namespaces     map[string]bool
		expectedOutput bool
	}{
		{
			namespace: "in-block-list",
			namespaces: map[string]bool{
				"kube-system":     true,
				"kube-public":     true,
				"kube-node-lease": true,
				"falco":           true,
				"in-block-list":   true,
			},
			expectedOutput: false,
		},
		{
			namespace: "not-critical",
			namespaces: map[string]bool{
				"kube-system":     true,
				"kube-public":     true,
				"kube-node-lease": true,
				"falco":           true,
			},
			expectedOutput: true,
		},
	}
	for _, single := range data {
		t.Run("", func(single struct {
			namespace      string
			namespaces     map[string]bool
			expectedOutput bool
		}) func(t *testing.T) {
			return func(t *testing.T) {
				output := CheckBlockNamespace(single.namespace, single.namespaces)
				if output != single.expectedOutput {
					t.Errorf("Got: %v, expected %v, for namespace %v", output, single.expectedOutput, single.namespace)
				}
			}
		}(single))
	}
}
