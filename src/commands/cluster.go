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

