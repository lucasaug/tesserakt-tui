package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type panelPosition string

const (
    Navigation panelPosition = "Navigation"
    Header     panelPosition = "Header"
    Main       panelPosition = "Main"
)

type direction string

const (
    up    direction = "up"
    left  direction = "left"
    down  direction = "down"
    right direction = "right"
)

// There's gotta be a better way to do this -.-'
var nextPanel map[panelPosition](map[direction]panelPosition) = map[panelPosition](map[direction]panelPosition){
    Navigation: map[direction]panelPosition {
        up: Navigation,
        left: Navigation,
        down: Navigation,
        right: Main,
    },
    Main: map[direction]panelPosition {
        up: Main,
        left: Navigation,
        down: Main,
        right: Main,
    },
}

type mainModel struct {
    width           int
    height          int

    navigation  resourcePicker
    mainContent resourceList

    currentPanel panelPosition
}

func InitialModel() mainModel {
    return mainModel{
        mainContent: InitialResourceListModel(),
        navigation: InitialResourcePickerModel(),
        currentPanel: Navigation,
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

        case "h", "left":
            m.currentPanel = nextPanel[m.currentPanel][left]

        case "l", "right":
            m.currentPanel = nextPanel[m.currentPanel][right]

        }
    }

    var cmd tea.Cmd

    if (m.currentPanel == Main) {
        m.mainContent.Focus()
        m.mainContent, cmd = m.mainContent.Update(msg)
    } else if (m.currentPanel == Navigation) {
        m.navigation.Focus()
        m.navigation, cmd = m.navigation.Update(msg)
    }

     m.mainContent.SetResource(Resources[m.navigation.resourceIndex])

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
            m.navigation.View(),
            m.mainContent.View(),
        ),
    )
}
