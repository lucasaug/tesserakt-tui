package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StatusBarUpdateMsg struct {
    ConnectionStatus string
    ClusterName string
    NodeData string
    Status string
}

type TickMsg time.Time

func Tick(tickInterval time.Duration) tea.Cmd {
    return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

