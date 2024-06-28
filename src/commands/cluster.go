package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasaug/tesserakt-tui/src/adapters"
	"github.com/lucasaug/tesserakt-tui/src/core"
	"github.com/lucasaug/tesserakt-tui/src/k8s"
	"k8s.io/client-go/kubernetes"
)

type K8sClientMsg struct { Clientset *kubernetes.Clientset }

func GetKubernetesClientCmd() tea.Msg {
    clientset := k8s.GetClientSet()
    return K8sClientMsg{clientset}
}

type ResourceDetailsMsg struct { Value core.ResourceSelector }

func ResourceDetails(
    clientset kubernetes.Clientset,
    resourceType core.Resource,
    name, namespace string,
) tea.Cmd {
    return func() tea.Msg {
        var data string
        if (resourceType == core.Pod) {
            data = adapters.GetPodJson(&clientset, name, namespace)
        } else if (resourceType == core.Deployment) {
            data = adapters.GetDeploymentJson(&clientset, name, namespace)
        } else if (resourceType == core.Ingress) {
            data = adapters.GetIngressJson(&clientset, name, namespace)
        }

        return ResourceDetailsMsg{
            Value: core.ResourceSelector{
                Name: name,
                Namespace: namespace,
                ResourceType: resourceType,
                Editable: false,
                Data: data,
            },
        }
    }
}


