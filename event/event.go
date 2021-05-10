package event

import (
	"encoding/json"
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

func CheckNamespace(namespace string, criticalNamespaces map[string]bool) bool {
	if criticalNamespaces[namespace] {
		return false
	}

	return true
}
