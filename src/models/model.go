package models

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/core"
	"github.com/lucasaug/tesserakt-tui/src/commands"
	"github.com/lucasaug/tesserakt-tui/src/k8s"
)

const TICK_INTERVAL = time.Millisecond * 500

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

var nextPanel map[panelPosition](map[direction]panelPosition) =
    map[panelPosition](map[direction]panelPosition){
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
    mainContent resourceView

    statusBar    statusbar.Model

    currentPanel panelPosition
}

func InitialModel() mainModel {
    mainContent := InitialResourceViewModel()
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


func (m mainModel) updateStatus() tea.Msg {
    nodes := k8s.GetNodes(m.clientset)

    nodeText := fmt.Sprint(len(nodes), " nodes")
    if len(nodes) == 1 {
        nodeText = "1 node"
    }

    return commands.StatusBarUpdateMsg{
        ConnectionStatus: "Connected",
        ClusterName: "cluster name",
        NodeData: nodeText,
        Status: "UP",
    }
}

func (m mainModel) Init() tea.Cmd {
    mainContentCmd := m.mainContent.Init()

    return tea.Sequence(
        commands.GetKubernetesClientCmd,
        mainContentCmd,
        commands.Tick(TICK_INTERVAL),
    )
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        navigationWidth := msg.Width / 8

        mainWidth := msg.Width - navigationWidth - 4

        componentHeight := msg.Height - m.statusBar.Height - 4

        m.navigation.SetSize(navigationWidth, componentHeight)
        m.mainContent.SetSize(mainWidth, componentHeight)
        m.statusBar.SetSize(msg.Width)

        var mainCmd, navigationCmd tea.Cmd
        m.mainContent, mainCmd = m.mainContent.Update(msg)
        m.navigation, navigationCmd = m.navigation.Update(msg)

        cmd = tea.Batch(mainCmd, navigationCmd)

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
            m.mainContent.SetHighlight(true)
            m.mainContent.Focus()
            m.navigation.SetHighlight(false)
            m.navigation.Blur()

            m.mainContent, cmd = m.mainContent.Update(msg)

        } else if (m.currentPanel == Navigation) {
            m.mainContent.SetHighlight(false)
            m.mainContent.Blur()
            m.navigation.SetHighlight(true)
            m.navigation.Focus()

            var navigationCmd, mainCmd tea.Cmd
            m.navigation, navigationCmd = m.navigation.Update(msg)
            m.mainContent, mainCmd = m.mainContent.Update(
                commands.ResourceChangeMsg{
                    NewResource: core.Resources[m.navigation.resourceIndex],
                },
            )

            cmd = tea.Batch(navigationCmd, mainCmd)
        }

    case commands.K8sClientMsg:
        m.clientset = msg.Clientset

        var originalMsgCmd, resourceChangeCmd tea.Cmd
        m.mainContent, originalMsgCmd = m.mainContent.Update(msg)
        m.mainContent, resourceChangeCmd = m.mainContent.Update(
            commands.ResourceChangeMsg{
                NewResource: core.Resources[m.navigation.resourceIndex],
            },
        )

        cmd = tea.Batch(
            m.updateStatus,
            resourceChangeCmd,
            originalMsgCmd,
        )

    case commands.TickMsg:
        var mainCmd tea.Cmd
        m.mainContent, mainCmd = m.mainContent.Update(msg)
        cmd = tea.Batch(
            m.updateStatus,
            mainCmd,
            commands.Tick(TICK_INTERVAL),
        )

    case commands.StatusBarUpdateMsg:
        m.statusBar.SetContent(
            msg.ConnectionStatus,
            msg.ClusterName,
            msg.NodeData,
            msg.Status,
        )

    default:
        m.mainContent, cmd = m.mainContent.Update(msg)

    }

    return m, cmd
}

func (m mainModel) View() string {
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
