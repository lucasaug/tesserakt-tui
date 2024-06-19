package k8s

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/networking/v1"
)


func GetIngresses() []v1.Ingress {
    clientset := GetClientSet()

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

