package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	adapters "github.com/lucasaug/tesserakt-tui/src/adapters"
)

type resourceList struct {
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

    itemListing := adapters.GetPodTable()

    return resourceList{
        table: itemListing,
        style: tableStyle,
        highlightedStyle: highlightedStyle,
    }
}

func (r resourceList) Init() tea.Cmd {
    return nil
}

func (r *resourceList) SetResource(res Resource) {
    if r.resourceType == res { return }

    r.resourceType = res
    if res == Pod {
        r.table = adapters.GetPodTable()
    } else if res == Deployment {
        r.table = adapters.GetDeploymentTable()
    } else if res == Ingress {
        r.table = adapters.GetIngressTable()
    } else {
        r.table = table.New()
    }

}

func (r resourceList) Update(msg tea.Msg) (resourceList, tea.Cmd) {
    var cmd tea.Cmd
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
