package event

import (
	"encoding/json"
	"strings"
	"time"
)

// Alert falco data structure
type Alert struct {
	Output       string    `json:"output"`
	Priority     string    `json:"priority"`
	Rule         string    `json:"rule"`
	Time         time.Time `json:"time"`
	OutputFields struct {
		ContainerID              string      `json:"container.id"`
		ContainerImageRepository interface{} `json:"container.image.repository"`
		ContainerImageTag        interface{} `json:"container.image.tag"`
		EvtTime                  int64       `json:"evt.time"`
		FdName                   string      `json:"fd.name"`
		K8SNsName                string      `json:"k8s.ns.name"`
		K8SPodName               string      `json:"k8s.pod.name"`
		ProcCmdline              string      `json:"proc.cmdline"`
	} `json:"output_fields"`
}

func Read(data []byte) (Alert, error) {
	var falcoEvent Alert

	err := json.Unmarshal(data, &falcoEvent)
	if err != nil {
		return Alert{}, err
	}

	return falcoEvent, nil
}

// CheckAllowNamespace if the namespace is in the hash map it's okay to delete in
func CheckAllowNamespace(namespace string, allowNamespaces map[string]bool) bool {
	if allowNamespaces[namespace] {
		return true
	}
	return false
}

// CheckBlockNamespace if the namespace is in the hash map the namespace is blocked for deletion
func CheckBlockNamespace(namespace string, blockNamespaces map[string]bool) bool {
	if blockNamespaces[namespace] {
		return false
	}
	return true
}

// AddItemsToHashMap namespaceInput should be a list of namespaces with space between, ex: 'ns1 ns2 app1-*'
func AddItemsToHashMap(namespaceInput string, namespaces map[string]bool) map[string]bool {
	newNamespaces := strings.Fields(namespaceInput)
	for _, s := range newNamespaces {
		namespaces[s] = true
	}

	return namespaces
}
