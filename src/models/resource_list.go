package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)

type resourceList struct {
    table table.Model
    style lipgloss.Style
}

func InitialResourceListModel() resourceList {
    columns := []table.Column {
        { Title: "Name", Width: 10 },
        { Title: "Namespace", Width: 10 },
        { Title: "Num of containers", Width: 20 },
        { Title: "Conditions", Width: 20 },
    }

    rows := []table.Row {}
    for _, pod := range k8s.GetPods() {
        rows = append(rows, pod)
    }

    borderColor := lipgloss.Color("36")
    tableStyle := lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)

    itemListing := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
    )

    return resourceList{
        table: itemListing,
        style: tableStyle,
    }
}

func (r resourceList) Init() tea.Cmd {
    return nil
}

func (r resourceList) Update(msg tea.Msg) (resourceList, tea.Cmd) {
    var cmd tea.Cmd
    r.table, cmd = r.table.Update(msg)

    return r, cmd
}

func (r resourceList) View() string {
    return r.style.Render(r.table.View())
}

