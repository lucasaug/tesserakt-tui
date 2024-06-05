package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Resource string

const (
    Pod        Resource = "Pod"
    Deployment Resource = "Deployment"
    Ingress    Resource = "Ingress"
)

type resourcePicker struct {
    currentResource Resource
    table           table.Model
    style           lipgloss.Style
}

func InitialResourcePickerModel() resourcePicker {
    resourceItems := []table.Row {
        []string{ string(Pod) },
        []string{ string(Deployment) },
        []string{ string(Ingress) },
    }

    borderColor := lipgloss.Color("36")

    listStyle := lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(20)

    itemListing := table.New(
        table.WithColumns([]table.Column {{ Title: "Resources", Width: 10 }}),
        table.WithRows(resourceItems),
        table.WithFocused(true),
    )

    return resourcePicker{
        currentResource: "Pod",
        table: itemListing,
        style: listStyle,
    }
}

func (r resourcePicker) Init() tea.Cmd {
    r.table.Focus()
    return nil
}

func (r resourcePicker) Update(msg tea.Msg) (resourcePicker, tea.Cmd) {
    var cmd tea.Cmd
    r.table, cmd = r.table.Update(msg)

    return r, cmd
}

func (r resourcePicker) View() string {
    return r.style.Render(r.table.View())
}

