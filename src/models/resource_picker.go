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

var Resources = [...]Resource{
    Pod,
    Deployment,
    Ingress,
}

type resourcePicker struct {
    resourceIndex int

    table         table.Model
    style         lipgloss.Style

    width int
    height int

    highlighted bool
}

func InitialResourcePickerModel() resourcePicker {
    resourceItems := []table.Row {
        []string{ string(Pod) },
        []string{ string(Deployment) },
        []string{ string(Ingress) },
    }
    resourceHeader := []table.Column {
        { Title: "Resources", Width: 10 },
    }

    listStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder())

    itemListing := table.New(
        table.WithColumns(resourceHeader),
        table.WithRows(resourceItems),
        table.WithFocused(true),
    )

    return resourcePicker{
        resourceIndex: 0,
        table: itemListing,
        style: listStyle,
    }
}

func (r resourcePicker) Init() tea.Cmd {
    return nil
}

func (r resourcePicker) Update(msg tea.Msg) (resourcePicker, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {

        case "k", "up":
            if (r.resourceIndex > 0) {
                r.resourceIndex--
            }

        case "j", "down":
            if (r.resourceIndex < len(Resources) - 1) {
                r.resourceIndex++
            }

        }
    }

    var cmd tea.Cmd
    r.table, cmd = r.table.Update(msg)

    return r, cmd
}

func (r resourcePicker) View() string {
    if (r.highlighted) {
        r.style.BorderForeground(highlightColor)
    } else {
        r.style.BorderForeground(borderColor)
    }

    r.table.SetWidth(r.width)
    r.table.SetHeight(r.height)

    return r.style.Render(r.table.View())
}

func (r resourcePicker) Focus() {
    r.table.Focus()
}

func (r resourcePicker) Blur() {
    r.table.Blur()
}

func (r *resourcePicker) SetSize(width int, height int) {
    r.width = width
    r.height = height
}

func (r *resourcePicker) SetHighlight(highlighted bool) {
    r.highlighted = highlighted
}
