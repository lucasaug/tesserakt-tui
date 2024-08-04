package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobResource v1.Job

func (pr JobResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
    }
}

func (pr JobResource) ResourceName() string {
    return pr.Name
}

func (pr JobResource) ResourceNamespace() string {
    return pr.Namespace
}

type JobHandler struct {
    list []JobResource
}

func (ph JobHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        BatchV1().
        Jobs(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    ph.list = []JobResource{}
    for _, item := range pods.Items {
        result = append(result, JobResource(item))
        ph.list = append(ph.list, JobResource(item))
    }

    return result, nil
}

func (_ JobHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
    }
}

