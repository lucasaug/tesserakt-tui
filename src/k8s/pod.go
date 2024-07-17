package k8s

import (
    "context"
    "fmt"

    "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type PodResource v1.Pod

func (pr PodResource) Values() []string {
    return []string{
        pr.Name,
        pr.Namespace,
        fmt.Sprint(len(pr.Spec.Containers)),
        string(pr.Status.Phase),
    }
}

type PodHandler struct {}

func (PodHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    pods, err := clientset.
        CoreV1().
        Pods(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    for _, item := range pods.Items {
        result = append(result, PodResource(item))
    }

    return result, nil
}

