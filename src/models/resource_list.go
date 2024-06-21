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

    table table.Model
    style lipgloss.Style
    highlightedStyle lipgloss.Style

    width int
    height int

    resourceType Resource

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

const TICK_INTERVAL = time.Second

type ResourceUpdateMsg struct { table table.Model }

func (r resourceList) createScheduleResourceUpdateCmd() tea.Cmd {
    return tea.Tick(TICK_INTERVAL, func(t time.Time) tea.Msg {
        return r.UpdateResourceCmd()
    })
}

func (r resourceList) UpdateResourceCmd() tea.Msg {
    var result table.Model

    if r.resourceType == Pod {
        result = adapters.GetPodTable(r.clientset)
    } else if r.resourceType == Deployment {
        result = adapters.GetDeploymentTable(r.clientset)
    } else if r.resourceType == Ingress {
        result = adapters.GetIngressTable(r.clientset)
    } else {
        result = table.New()
    }

    return ResourceUpdateMsg{table: result}
}

func (r resourceList) Init() tea.Cmd {
    return nil
}

func (r resourceList) Update(msg tea.Msg) (resourceList, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case k8sClientMsg:
        r.clientset = msg.clientset
        return r, r.UpdateResourceCmd

    case ResourceUpdateMsg:
        r.table = msg.table
        return r, r.UpdateResourceCmd

    }

    r.table, cmd = r.table.Update(msg)
    return r, cmd
}

func (r resourceList) View() string {
    r.table.SetWidth(r.width)
    r.table.SetHeight(r.height)

    if (r.highlighted) {
        return r.highlightedStyle.Render(r.table.View())
    }

    return r.style.Render(r.table.View())
}

func (r resourceList) Focus() {
    r.table.Focus()
}

func (r resourceList) Blur() {
    r.table.Blur()
}

func (r *resourceList) SetSize(width int, height int) {
    r.width = width
    r.height = height
}

func (r *resourceList) SetHighlight(highlighted bool) {
    r.highlighted = highlighted
}
