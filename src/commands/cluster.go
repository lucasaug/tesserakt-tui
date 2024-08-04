package commands

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/lucasaug/tesserakt-tui/src/k8s"
    "k8s.io/client-go/kubernetes"
)

type K8sClientMsg struct { Clientset *kubernetes.Clientset }
type K8sClusterNameMsg struct { Name string }

func GetKubernetesClientCmd() tea.Msg {
    return K8sClientMsg{k8s.GetClientSet()}
}

func GetClusterNameCmd() tea.Msg {
    return K8sClusterNameMsg{k8s.GetClusterName()}
}

