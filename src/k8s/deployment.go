package k8s

import (
	"context"
	"fmt"
	"os"

	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func GetDeployments(clientset *kubernetes.Clientset) []v1.Deployment {
    deployments, err := clientset.
        AppsV1().
        Deployments("").
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        fmt.Printf("error getting deployments: %v\n", err)
        os.Exit(1)
    }

    return deployments.Items
}

