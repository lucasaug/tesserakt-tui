package k8s

import (
	"context"
	"fmt"
	"os"

	v1Corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func GetPods(clientset *kubernetes.Clientset) []v1Corev1.Pod {
    pods, err := clientset.
        CoreV1().
        Pods("").
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        fmt.Printf("error getting pods: %v\n", err)
        os.Exit(1)
    }

    return pods.Items
}
