package k8s

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/apps/v1"
)


func GetDeployments() []v1.Deployment {
    clientset := GetClientSet()

    deployments, err := clientset.
        AppsV1().
        Deployments("default").
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        fmt.Printf("error getting deployments: %v\n", err)
        os.Exit(1)
    }

    return deployments.Items
}

