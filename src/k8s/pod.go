package k8s

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1Corev1 "k8s.io/api/core/v1"
)


func GetPods() []v1Corev1.Pod {
    clientset := GetClientSet()

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

