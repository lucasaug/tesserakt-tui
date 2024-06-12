package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
        "github.com/mistakenelf/teacup/statusbar"
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
    statusBar    statusbar.Model
}

func InitialModel() mainModel {
    sb := statusbar.New(
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
            Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
        },
    )

    return mainModel{
        mainContent: InitialResourceListModel(),
        navigation: InitialResourcePickerModel(),
        currentPanel: Navigation,
        statusBar: sb,
    }
}

func (m mainModel) Init() tea.Cmd {
    return nil;
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        navigationWidth := msg.Width / 5
        if navigationWidth > 300 {
            navigationWidth = 300
        }

        m.mainContent.SetSize(msg.Width - navigationWidth, msg.Height)
        m.navigation.SetSize(navigationWidth, msg.Height)

        m.statusBar.SetSize(msg.Width)
        m.statusBar.SetContent("Connected", "my-cluster-prd", "192.168.0.1", "UP")

    case tea.KeyMsg:
        switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit

        case "h", "left":
            m.currentPanel = nextPanel[m.currentPanel][left]

        case "l", "right":
            m.currentPanel = nextPanel[m.currentPanel][right]

        }

        if (m.currentPanel == Main) {
            m.mainContent.Focus()
            m.mainContent, cmd = m.mainContent.Update(msg)
        } else if (m.currentPanel == Navigation) {
            m.navigation.Focus()
            m.navigation, cmd = m.navigation.Update(msg)
        }

        m.mainContent.SetResource(Resources[m.navigation.resourceIndex])
    }

    return m, cmd
}

func (m mainModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Bottom,
        lipgloss.JoinHorizontal(
            lipgloss.Top,
            m.navigation.View(),
            m.mainContent.View(),
        ),
        m.statusBar.View(),
    )
}
