package k8s

import (
	"slices"

	"k8s.io/client-go/kubernetes"
)

type ResourceType string

const (
    Pod        ResourceType = "Pod"
    Deployment ResourceType = "Deployment"
    Ingress    ResourceType = "Ingress"
)

var Resources = [...]ResourceType{
    Pod,
    Deployment,
    Ingress,
}

type ResourceSelector struct {
    Name string
    Namespace string
    Editable bool
    Data string
}

type ResourceInstance interface {
    Values() []string
    ResourceName() string
}

type NamespacedResource interface {
    ResourceInstance
    ResourceNamespace() string
}

type ResourceHandler interface {
    List(*kubernetes.Clientset, string) ([]ResourceInstance, error)
}

func GetResource(
    clientset *kubernetes.Clientset,
    rh ResourceHandler,
    name, namespace string,
) ResourceInstance {
    // TODO handle err
    list, _ := rh.List(clientset, namespace)

    idx := slices.IndexFunc(list, func(r ResourceInstance) bool {
        resource, ok := r.(NamespacedResource)
        if !ok {
            return false
        }

        return resource.ResourceName() == name &&
            resource.ResourceNamespace() == namespace
    })

    return list[idx]
}
