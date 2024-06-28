package commands

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/adapters"
	"github.com/lucasaug/tesserakt-tui/src/core"
)

type LoadTablesMsg struct { Tables map[core.Resource]*table.Model }
func CreateLoadTables(clientset *kubernetes.Clientset) tea.Cmd {
    return func() tea.Msg {
        podTable := adapters.GetPodTable(clientset)
        deploymentTable := adapters.GetDeploymentTable(clientset)
        ingressTable := adapters.GetIngressTable(clientset)

        tables := map[core.Resource]*table.Model {
            core.Pod: &podTable,
            core.Deployment: &deploymentTable,
            core.Ingress: &ingressTable,
        }

        return LoadTablesMsg{Tables: tables}
    }
}

type ResourceChangeMsg struct { NewResource core.Resource }

type RefreshResourceListMsg struct {
    Rows []table.Row
    Resource core.Resource
}
type EmptyMsg struct {}

func RefreshResourceList(
    clientset *kubernetes.Clientset,
    resourceType *core.Resource,
) tea.Cmd {
    return func() tea.Msg {
	if (clientset == nil) { return EmptyMsg{} }

	if *resourceType == core.Pod {
	    return RefreshResourceListMsg{
		Rows: adapters.GetPodRows(clientset),
		Resource: *resourceType,
	    }
	} else if *resourceType == core.Deployment {
	    return RefreshResourceListMsg{
		Rows: adapters.GetDeploymentRows(clientset),
		Resource: *resourceType,
	    }
	} else if *resourceType == core.Ingress {
	    return RefreshResourceListMsg{
		Rows: adapters.GetIngressRows(clientset),
		Resource: *resourceType,
	    }
	}

	return RefreshResourceListMsg{
	    Rows: []table.Row{},
	    Resource: *resourceType,
	}
    }
}

