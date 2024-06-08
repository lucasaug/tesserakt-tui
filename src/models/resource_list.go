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

    resourceType Resource
}

func InitialResourceListModel() resourceList {
    borderColor := lipgloss.Color("36")
    tableStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder()).
        Padding(1).
        Width(80)

    itemListing := adapters.GetPodTable()
    return resourceList{
        table: itemListing,
        style: tableStyle,
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
    return r.style.Render(r.table.View())
}

