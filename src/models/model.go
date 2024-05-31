package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

    k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)

type Resource string

const (
    Pod        Resource = "Pod"
    Deployment Resource = "Deployment"
    Ingress    Resource = "Ingress"
)

type model struct {
    currentResource Resource
    currentIndex    int
    resourceList    table.Model
    table           table.Model
    width           int
    height          int
    tableStyle      lipgloss.Style
    listStyle       lipgloss.Style
}

func InitialModel() model {
    resourceItems := []table.Row {
        []string{ string(Pod) },
        []string{ string(Deployment) },
        []string{ string(Ingress) },
    }

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
    listStyle := lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(20)

    itemListing := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(7),
    )

    return model{
        currentResource: "Pod",
        currentIndex: 0,
        table: itemListing,
        resourceList: table.New(
            table.WithColumns([]table.Column {{ Title: "Resources", Width: 10 }}),
            table.WithRows(resourceItems),
            table.WithFocused(true),
            table.WithHeight(7),
        ),
        tableStyle: tableStyle,
        listStyle: listStyle,
    }
}

func (m model) Init() tea.Cmd {
    return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case tea.KeyMsg:

        switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit

        // case "up", "k":
        //     if m.currentIndex < len(m.names) {
        //         m.currentIndex++
        //     }
        //
        // case "down", "j":
        //     if m.currentIndex > 0 {
        //         m.currentIndex--
        //     }
        //
        }
    }

    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)

    return m, cmd
}

func (m model) View() string {
    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        lipgloss.JoinHorizontal(
            lipgloss.Top,
            m.listStyle.Render(m.resourceList.View()),
            m.tableStyle.Render(m.table.View()),
        ),
    )
}
