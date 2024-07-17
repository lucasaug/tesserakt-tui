package commands

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/lucasaug/tesserakt-tui/src/k8s"
    "k8s.io/client-go/kubernetes"
)

type K8sClientMsg struct { Clientset *kubernetes.Clientset }

func GetKubernetesClientCmd() tea.Msg {
    clientset := k8s.GetClientSet()
    return K8sClientMsg{clientset}
}

type ResourceDetailsMsg struct { Value k8s.ResourceSelector }

func ResourceDetails(
    clientset kubernetes.Clientset,
    resourceType k8s.ResourceType,
    name, namespace string,
) tea.Cmd {
    return func() tea.Msg {
        // var data string
        // if (resourceType == k8s.Pod) {
        //     data = adapters.GetPodJson(&clientset, name, namespace)
        // } else if (resourceType == k8s.Deployment) {
        //     data = adapters.GetDeploymentJson(&clientset, name, namespace)
        // } else if (resourceType == k8s.Ingress) {
        //     data = adapters.GetIngressJson(&clientset, name, namespace)
        // }
        //
        // return ResourceDetailsMsg{
        //     Value: k8s.ResourceSelector{
        //         Name: name,
        //         Namespace: namespace,
        //         Editable: false,
        //         Data: data,
        //     },
        // }
        return EmptyMsg{}
    }
}


