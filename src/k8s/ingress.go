package k8s

import (
    "context"

    v1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type IngressResource v1.Ingress

func (ir IngressResource) Values() []string {
    return []string{
        ir.Name,
        ir.Namespace,
    }
}

type IngressHandler struct {}

func (IngressHandler) List(
    clientset *kubernetes.Clientset,
    namespace string,
) ([]ResourceInstance, error) {
    ingresses, err := clientset.
        NetworkingV1().
        Ingresses(namespace).
        List(context.Background(), metav1.ListOptions{})

    if err != nil {
        return []ResourceInstance{}, err
    }

    result := []ResourceInstance{}
    for _, item := range ingresses.Items {
        result = append(result, IngressResource(item))
    }

    return result, nil
}

