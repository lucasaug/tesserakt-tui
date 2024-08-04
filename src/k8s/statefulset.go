package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StatefulSetResource v1.StatefulSet

func (pr StatefulSetResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
    }
}

func (pr StatefulSetResource) ResourceName() string {
    return pr.Name
}

func (pr StatefulSetResource) ResourceNamespace() string {
    return pr.Namespace
}

type StatefulSetHandler struct {
    list []StatefulSetResource
}

func (ph StatefulSetHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        AppsV1().
        StatefulSets(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    ph.list = []StatefulSetResource{}
    for _, item := range pods.Items {
        result = append(result, StatefulSetResource(item))
        ph.list = append(ph.list, StatefulSetResource(item))
    }

    return result, nil
}

func (_ StatefulSetHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
    }
}

