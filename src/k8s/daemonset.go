package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DaemonSetResource v1.DaemonSet

func (pr DaemonSetResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
    }
}

func (pr DaemonSetResource) ResourceName() string {
    return pr.Name
}

func (pr DaemonSetResource) ResourceNamespace() string {
    return pr.Namespace
}

type DaemonSetHandler struct {
    list []DaemonSetResource
}

func (ph DaemonSetHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        AppsV1().
        DaemonSets(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    ph.list = []DaemonSetResource{}
    for _, item := range pods.Items {
        result = append(result, DaemonSetResource(item))
        ph.list = append(ph.list, DaemonSetResource(item))
    }

    return result, nil
}

func (_ DaemonSetHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
    }
}

