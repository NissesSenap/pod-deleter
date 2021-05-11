package poddelete

import (
	"context"
	"fmt"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckPodAnnotation(t *testing.T) {
	// testCases
	data := []struct {
		expectedBool   bool
		inputNamespace string
		inputName      string
		err            error
	}{
		{
			expectedBool:   false,
			inputNamespace: "default",
			inputName:      "missPod",
			err:            fmt.Errorf("Unable to get pod: missPod in namespace: default, due to error: pods \"missPod\" not found"),
		},
		{
			expectedBool:   true,
			inputNamespace: "default",
			inputName:      "truePod",
			err:            nil,
		},
		{
			expectedBool:   false,
			inputNamespace: "default",
			inputName:      "falsePod1",
			err:            nil,
		},
		{
			expectedBool:   false,
			inputNamespace: "default",
			inputName:      "falsePod2",
			err:            nil,
		},
		{
			expectedBool:   false,
			inputNamespace: "default",
			inputName:      "falsePod3",
			err:            nil,
		},
	}

	// insert fake client data
	clientSet := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        "truePod",
			Namespace:   "default",
			Annotations: map[string]string{falcoAnnotation: "true"},
		},
	}, &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        "falsePod1",
			Namespace:   "default",
			Annotations: map[string]string{falcoAnnotation: "false"},
		},
	}, &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:        "falsePod2",
			Namespace:   "default",
			Annotations: map[string]string{},
		},
	}, &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "falsePod3",
			Namespace: "default",
		},
	})

	testClient := &PodDeleter{kubeClient: clientSet}

	ctx := context.Background()

	for _, single := range data {
		t.Run("", func(single struct {
			expectedBool   bool
			inputNamespace string
			inputName      string
			err            error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				output, err := testClient.CheckPodAnnotation(ctx, single.inputNamespace, single.inputName)
				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				} else {
					if output != single.expectedBool {
						t.Errorf("Got output: %v, expected %v", output, single.expectedBool)
					}
				}
			}
		}(single))
	}
}

func TestDeletePod(t *testing.T) {
	// testCases
	data := []struct {
		inputNamespace string
		inputName      string
		err            error
	}{
		{
			inputNamespace: "default",
			inputName:      "missPod",
			err:            fmt.Errorf("pods \"missPod\" not found"),
		},
		{
			inputNamespace: "default",
			inputName:      "annotatedPod",
			err:            nil,
		},
		{
			inputNamespace: "default",
			inputName:      "nonAnnotatedPod",
			err:            nil,
		},
	}

	// insert fake client data
	clientSet := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "annotatedPod",
			Namespace: "default",
			// The check is done earlier, if the pod magnically have the time to get a annotation
			// after the annotation check the pod will still be deleted.
			Annotations: map[string]string{falcoAnnotation: "true"},
		},
	}, &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "nonAnnotatedPod",
			Namespace: "default",
		},
	})

	testClient := &PodDeleter{kubeClient: clientSet}

	ctx := context.Background()

	for _, single := range data {
		t.Run("", func(single struct {
			inputNamespace string
			inputName      string
			err            error
		}) func(t *testing.T) {
			return func(t *testing.T) {
				err := testClient.DeletePod(ctx, single.inputNamespace, single.inputName)
				if err != nil {
					if single.err == nil {
						t.Fatalf(err.Error())
					}
					if !strings.EqualFold(single.err.Error(), err.Error()) {
						t.Fatalf("expected err: %s got err: %s", single.err, err)
					}
				}
			}
		}(single))
	}
}
