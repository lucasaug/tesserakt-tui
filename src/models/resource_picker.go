package models

import (
    "github.com/charmbracelet/bubbles/table"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/lucasaug/tesserakt-tui/src/k8s"
)

type resourcePicker struct {
    width int
    height int

    style            lipgloss.Style
    highlightedStyle lipgloss.Style
    table            table.Model

    resourceIndex int
    highlighted bool
}

func InitialResourcePickerModel() resourcePicker {
    resourceItems := []table.Row {
        []string{ string(k8s.Pod) },
        []string{ string(k8s.Deployment) },
        []string{ string(k8s.Ingress) },
    }
    resourceHeader := []table.Column {
        { Title: "Resources", Width: 10 },
    }

    listStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder())

    highlightedStyle := lipgloss.NewStyle().
        BorderForeground(highlightColor).
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
        highlightedStyle: highlightedStyle,
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
            if (r.resourceIndex < len(k8s.Resources) - 1) {
                r.resourceIndex++
            }

        }
    }

    var cmd tea.Cmd
    r.table, cmd = r.table.Update(msg)

    return r, cmd
}

func (r resourcePicker) View() string {
    r.table.SetWidth(r.width)
    r.table.SetHeight(r.height)

    if (r.highlighted) {
        return r.highlightedStyle.Render(r.table.View())
    }
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
