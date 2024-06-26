package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNodes(clientset *kubernetes.Clientset) []v1.Node {
    nodes, _ := clientset.CoreV1().Nodes().List(
        context.TODO(),
        metav1.ListOptions{},
    )

    return nodes.Items
}

