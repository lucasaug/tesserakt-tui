package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/k8s"
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
    clientset *kubernetes.Clientset

    width           int
    height          int

    navigation  resourcePicker
    mainContent resourceList

    currentPanel panelPosition
    statusBar    statusbar.Model
}

func createStatusBar() statusbar.Model {
    sb := statusbar.New(
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Dark: "15", Light: "15"},
            Background: lipgloss.AdaptiveColor{Light: "13", Dark: "13"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "238", Dark: "238"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "171", Dark: "171"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "93", Dark: "93"},
        },
    )

    sb.SetContent("Connected", "my-cluster-prd", "192.168.0.1", "UP")

    return sb
}

func InitialModel() mainModel {
    mainContent := InitialResourceListModel()
    navigation := InitialResourcePickerModel()

    mainContent.SetHighlight(true)
    navigation.SetHighlight(false)

    return mainModel{
        mainContent: mainContent,
        navigation: navigation,
        currentPanel: Navigation,
        statusBar: createStatusBar(),
    }
}

func getKubernetesClient() tea.Msg {
    clientset := k8s.GetClientSet()
    return k8sClientMsg{clientset}
}

type k8sClientMsg struct { clientset *kubernetes.Clientset }

func (m mainModel) Init() tea.Cmd {
    mainContentCmd := m.mainContent.Init()
    return tea.Sequence(getKubernetesClient, mainContentCmd)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {

    case k8sClientMsg:
        m.clientset = msg.clientset
        m.mainContent, cmd = m.mainContent.Update(msg)
        m.mainContent.SetResource(Resources[m.navigation.resourceIndex])

        return m, tea.Batch(m.mainContent.refreshList, cmd)

    case RefreshListMsg, TickMsg:
        m.mainContent, cmd = m.mainContent.Update(msg)
        return m, cmd

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        navigationWidth := msg.Width / 8

        mainWidth := msg.Width - navigationWidth - 4

        componentHeight := msg.Height - m.statusBar.Height - 4

        m.navigation.SetSize(navigationWidth, componentHeight)
        m.mainContent.SetSize(mainWidth, componentHeight)

        m.statusBar.SetSize(msg.Width)

        m.mainContent, cmd = m.mainContent.Update(msg)
        m.navigation, cmd = m.navigation.Update(msg)

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
            m.navigation.Blur()
            m.mainContent, cmd = m.mainContent.Update(msg)
        } else if (m.currentPanel == Navigation) {
            m.mainContent.Blur()
            m.navigation.Focus()

            previousResourceIndex := m.navigation.resourceIndex
            m.navigation, cmd = m.navigation.Update(msg)

            if previousResourceIndex != m.navigation.resourceIndex {
                m.mainContent.SetResource(
                    Resources[m.navigation.resourceIndex],
                )
                return m, tea.Sequence(cmd, m.mainContent.refreshList)
            }
        }

    }

    return m, cmd
}

func (m mainModel) View() string {
    if (m.currentPanel == Main) {
        m.mainContent.SetHighlight(true)
        m.navigation.SetHighlight(false)
    } else if (m.currentPanel == Navigation) {
        m.mainContent.SetHighlight(false)
        m.navigation.SetHighlight(true)
    }

    return lipgloss.JoinVertical(
        lipgloss.Top,
        lipgloss.JoinHorizontal(
            lipgloss.Top,
            m.navigation.View(),
            m.mainContent.View(),
        ),
        m.statusBar.View(),
    )
}
