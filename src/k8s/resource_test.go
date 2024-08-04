package k8s

import (
	"errors"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	"testing"
)

type MockResource struct {
    name, namespace string
}

func (mr MockResource) Values() []string {
    return []string{}
}

func (mr MockResource) ResourceName() string {
    return mr.name
}

func (mr MockResource) ResourceNamespace() string {
    return mr.namespace
}

type MockHandler struct {
    shouldFail bool
}

func (_ MockHandler) Columns() []table.Column {
    return []table.Column{}
}

func (mh MockHandler) List(
    _ *kubernetes.Clientset,
    _ string,
) ([]ResourceInstance, error) {
    if (mh.shouldFail) {
        return []ResourceInstance{}, errors.New("Failed")
    }

    return []ResourceInstance{
        MockResource{name: "test1", namespace: "test-namespace"},
        MockResource{name: "test2", namespace: "another-namespace"},
    }, nil
}

func TestGetResource(t *testing.T) {
    clientset := kubernetes.Clientset{}
    handler := MockHandler{shouldFail: false}
    
    result := GetResource(&clientset, handler, "test1", "test-namespace")
    if result.ResourceName() != "test1" {
        t.Fatalf("Failed getting resource")
    }

}

