package k8s

import (
    "context"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type NodeResource v1.Node

func (nr NodeResource) Values() []string {
    return []string{
	nr.Name,
    }
}

type NodeHandler struct {}

func (NodeHandler) List(
    clientset *kubernetes.Clientset,
) ([]NodeResource, error) {
    nodes, err := clientset.CoreV1().Nodes().List(
        context.TODO(),
        metav1.ListOptions{},
    )

    if err != nil {
	return []NodeResource{}, err
    }

    result := []NodeResource{}
    for _, item := range nodes.Items {
        result = append(result, NodeResource(item))
    }

    return result, nil
}

