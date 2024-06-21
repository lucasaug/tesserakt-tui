package k8s

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func GetIngresses(clientset *kubernetes.Clientset) []v1.Ingress {
    pods, err := clientset.
        NetworkingV1().
        Ingresses("").
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        fmt.Printf("error getting pods: %v\n", err)
        os.Exit(1)
    }

    return pods.Items
}

