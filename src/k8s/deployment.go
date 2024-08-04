package k8s

import (
	"context"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentResource v1.Deployment

func (dr DeploymentResource) Values() []string {
    return []string{
        dr.Name,
        dr.Namespace,
    }
}

func (dr DeploymentResource) ResourceName() string {
    return dr.Name
}

func (dr DeploymentResource) ResourceNamespace() string {
    return dr.Namespace
}

type DeploymentHandler struct {}

func (DeploymentHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    deployments, err := clientset.
        AppsV1().
        Deployments(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    for _, item := range deployments.Items {
        result = append(result, DeploymentResource(item))
    }

    return result, nil
}

func (_ DeploymentHandler) Columns() []table.Column {
    return []table.Column{
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    }
}
