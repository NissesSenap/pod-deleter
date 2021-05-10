package delete

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

/*
func MockSetupKubeClient(internal bool) (*PodDeleter, error) {
	clientSet := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "influxdb-v2",
			Namespace: "default",
		},
	})
	return &PodDeleter{kubeClient: clientSet}, nil
}
*/

func TestCheckPodAnnotation(t *testing.T) {
	clientSet := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "influxdb-v2",
			Namespace: "default",
		},
	})

	testClient := &PodDeleter{kubeClient: clientSet}

	//stuff, _ := MockSetupKubeClient(false)
	ctx := context.Background()

	answer, _ := testClient.CheckPodAnnotation(ctx, "default", "influxdb-v2")
	if answer {
		t.Fatalf("Expected false got: %v", answer)
	}
}
