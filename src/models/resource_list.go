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

const TICK_INTERVAL = time.Millisecond * 500

type TickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(TICK_INTERVAL, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

type RefreshListMsg struct { table table.Model }
type EmptyMsg struct {}

func (r resourceList) refreshList() tea.Msg {
    if (r.clientset == nil) { return EmptyMsg{} }

    if r.resourceType == Pod {
        return RefreshListMsg{
            table: adapters.GetPodTable(r.clientset),
        }
    } else if r.resourceType == Deployment {
        return RefreshListMsg{
            table: adapters.GetDeploymentTable(r.clientset),
        }
    } else if r.resourceType == Ingress {
        return RefreshListMsg{
            table: adapters.GetIngressTable(r.clientset),
        }
    }

    return RefreshListMsg{
        table: table.New(),
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

    case TickMsg:
        return r, tea.Batch(r.refreshList, tick())

    case RefreshListMsg:
        r.table = msg.table

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
