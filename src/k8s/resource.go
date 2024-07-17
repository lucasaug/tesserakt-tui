package k8s

import (
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
}

type ResourceHandler interface {
    List(*kubernetes.Clientset, string) ([]ResourceInstance, error)
}

