package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainModel struct {
    width           int
    height          int

    resourcePicker  resourcePicker
    resourceList    resourceList
}

func InitialModel() mainModel {
    return mainModel{
        resourceList: InitialResourceListModel(),
        resourcePicker: InitialResourcePickerModel(),
    }
}

func (m mainModel) Init() tea.Cmd {
    return nil;
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
    m.resourcePicker, cmd = m.resourcePicker.Update(msg)

    m.resourceList.SetResource(Resources[m.resourcePicker.resourceIndex])

    return m, cmd
}

func (m mainModel) View() string {
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
