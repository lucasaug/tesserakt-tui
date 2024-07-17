package commands

import (
    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "k8s.io/client-go/kubernetes"

    "github.com/lucasaug/tesserakt-tui/src/controllers"
    "github.com/lucasaug/tesserakt-tui/src/k8s"
)

type ResourceChangeMsg struct { NewResource k8s.ResourceType }

type RefreshResourceListMsg struct {
    Rows []table.Row
    Resource k8s.ResourceType
}
type EmptyMsg struct {}

func RefreshResourceList(
    clientset *kubernetes.Clientset,
    resourceType *k8s.ResourceType,
) tea.Cmd {
    return func() tea.Msg {
	if (clientset == nil) { return EmptyMsg{} }

	// TODO handle err
	rows, _ := controllers.GetRows(clientset, *resourceType)
	return RefreshResourceListMsg{ Rows: rows, Resource: *resourceType }
    }
}

