package k8s

import (
	"context"
	"fmt"
	"os"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func GetPods(clientset *kubernetes.Clientset) []v1.Pod {
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

func GetPod(
    clientset *kubernetes.Clientset,
    name, namespace string,
) *v1.Pod {
    pod, err := clientset.
        CoreV1().
        Pods(namespace).
	Get(context.Background(), name, metav1.GetOptions{})

    if err != nil {
        fmt.Printf("error getting pod: %v\n", err)
        os.Exit(1)
    }

    return pod
}

