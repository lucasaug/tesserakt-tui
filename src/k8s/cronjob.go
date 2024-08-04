package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CronJobResource v1.CronJob

func (pr CronJobResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
    }
}

func (pr CronJobResource) ResourceName() string {
    return pr.Name
}

func (pr CronJobResource) ResourceNamespace() string {
    return pr.Namespace
}

type CronJobHandler struct {
    list []CronJobResource
}

func (ph CronJobHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        BatchV1().
        CronJobs(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    ph.list = []CronJobResource{}
    for _, item := range pods.Items {
        result = append(result, CronJobResource(item))
        ph.list = append(ph.list, CronJobResource(item))
    }

    return result, nil
}

func (_ CronJobHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
    }
}

