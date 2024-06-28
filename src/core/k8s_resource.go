package core

type Resource string

const (
    Pod        Resource = "Pod"
    Deployment Resource = "Deployment"
    Ingress    Resource = "Ingress"
)

var Resources = [...]Resource{
    Pod,
    Deployment,
    Ingress,
}

type ResourceSelector struct {
    Name string
    Namespace string
    ResourceType Resource
    Editable bool
    Data string
}

