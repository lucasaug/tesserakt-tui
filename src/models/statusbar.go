package models

import (
    "github.com/charmbracelet/lipgloss"
    "github.com/mistakenelf/teacup/statusbar"
)

func createStatusBar() statusbar.Model {
    sb := statusbar.New(
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Dark: "15", Light: "15"},
            Background: lipgloss.AdaptiveColor{Light: "2", Dark: "2"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "238", Dark: "238"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "93", Dark: "93"},
        },
        statusbar.ColorConfig{
            Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "15"},
            Background: lipgloss.AdaptiveColor{Light: "93", Dark: "93"},
       },
    )

    sb.SetContent("Connected", "cluster name", "", "")

    return sb
}

