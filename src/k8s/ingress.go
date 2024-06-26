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
        fmt.Printf("error getting ingress: %v\n", err)
        os.Exit(1)
    }

    return pods.Items
}

func GetIngress(
    clientset *kubernetes.Clientset,
    name, namespace string,
) *v1.Ingress {
    ingress, err := clientset.
        NetworkingV1().
        Ingresses(namespace).
	Get(context.Background(), name, metav1.GetOptions{})

    if err != nil {
        fmt.Printf("error getting ingress: %v\n", err)
        os.Exit(1)
    }

    return ingress
}
