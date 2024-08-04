package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ReplicaSetResource v1.ReplicaSet

func (pr ReplicaSetResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
    }
}

func (pr ReplicaSetResource) ResourceName() string {
    return pr.Name
}

func (pr ReplicaSetResource) ResourceNamespace() string {
    return pr.Namespace
}

type ReplicaSetHandler struct {
    list []ReplicaSetResource
}

func (ph ReplicaSetHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        AppsV1().
        ReplicaSets(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    ph.list = []ReplicaSetResource{}
    for _, item := range pods.Items {
        result = append(result, ReplicaSetResource(item))
        ph.list = append(ph.list, ReplicaSetResource(item))
    }

    return result, nil
}

func (_ ReplicaSetHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
    }
}

