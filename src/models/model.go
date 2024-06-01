package models

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    width           int
    height          int
    resourcePicker  resourcePicker
    resourceList    resourceList
}

func InitialModel() model {
    return model{
        resourceList: InitialResourceListModel(),
        resourcePicker: InitialResourcePickerModel(),
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

        }
    }

    var cmd tea.Cmd
    m.resourceList, cmd = m.resourceList.Update(msg)

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
            m.resourcePicker.View(),
            m.resourceList.View(),
        ),
    )
}
