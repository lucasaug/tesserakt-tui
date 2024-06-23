package models

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/adapters"
)

type resourceList struct {
    clientset *kubernetes.Clientset

    style lipgloss.Style
    highlightedStyle lipgloss.Style

    width int
    height int

    resourceType Resource
    resourceTables map[Resource]*table.Model

    highlighted bool
}

func InitialResourceListModel() resourceList {
    tableStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder())
    highlightedStyle := lipgloss.NewStyle().
        BorderForeground(highlightColor).
        BorderStyle(lipgloss.NormalBorder())

    return resourceList{
        style: tableStyle,
        highlightedStyle: highlightedStyle,
    }
}

func (r *resourceList) SetResource(res Resource) {
    r.resourceType = res
}

const TICK_INTERVAL = time.Millisecond * 500

type TickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(TICK_INTERVAL, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

type RefreshRowsMsg struct {
    rows []table.Row
    resource Resource
}
type LoadTablesMsg struct { tables map[Resource]*table.Model }
type EmptyMsg struct {}

func (r resourceList) refreshList() tea.Msg {
    if (r.clientset == nil) { return EmptyMsg{} }

    if r.resourceType == Pod {
        return RefreshRowsMsg{
            rows: adapters.GetPodRows(r.clientset),
            resource: r.resourceType,
        }
    } else if r.resourceType == Deployment {
        return RefreshRowsMsg{
            rows: adapters.GetDeploymentRows(r.clientset),
            resource: r.resourceType,
        }
    } else if r.resourceType == Ingress {
        return RefreshRowsMsg{
            rows: adapters.GetIngressRows(r.clientset),
            resource: r.resourceType,
        }
    }

    return RefreshRowsMsg{
        rows: []table.Row{},
        resource: r.resourceType,
    }
}

func createLoadTables(clientset *kubernetes.Clientset) tea.Cmd {
    return func() tea.Msg {
        podTable := adapters.GetPodTable(clientset)
        deploymentTable := adapters.GetDeploymentTable(clientset)
        ingressTable := adapters.GetIngressTable(clientset)

        tables := map[Resource]*table.Model {
            Pod: &podTable,
            Deployment: &deploymentTable,
            Ingress: &ingressTable,
        }

        return LoadTablesMsg{tables: tables}
    }
}

func (r resourceList) Init() tea.Cmd {
    return tick()
}

func (r resourceList) Update(msg tea.Msg) (resourceList, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case k8sClientMsg:
        r.clientset = msg.clientset
        return r, createLoadTables(r.clientset)

    case TickMsg:
        return r, tea.Batch(r.refreshList, tick())

    case LoadTablesMsg:
        r.resourceTables = msg.tables

    case RefreshRowsMsg:
        if (len(r.resourceTables) != 0) {
            r.resourceTables[msg.resource].SetRows(msg.rows)
        }

    }

    if (len(r.resourceTables) != 0) {
        *r.resourceTables[r.resourceType], cmd = r.resourceTables[r.resourceType].Update(msg)
    }

    return r, cmd
}

func (r resourceList) View() string {
    if (len(r.resourceTables) == 0) { return "" }

    r.resourceTables[r.resourceType].SetWidth(r.width)
    r.resourceTables[r.resourceType].SetHeight(r.height)

    if (r.highlighted) {
        return r.highlightedStyle.Render(
            r.resourceTables[r.resourceType].View(),
        )
    }

    return r.style.Render(r.resourceTables[r.resourceType].View())
}

func (r resourceList) Focus() {
    r.resourceTables[r.resourceType].Focus()
}

func (r resourceList) Blur() {
    r.resourceTables[r.resourceType].Blur()
}

func (r *resourceList) SetSize(width int, height int) {
    r.width = width
    r.height = height
}

func (r *resourceList) SetHighlight(highlighted bool) {
    r.highlighted = highlighted
}
