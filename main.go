package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"

    model "github.com/lucasaug/tesserakt-tui/src/models"
)

func main() {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        fmt.Println("Error creating log file:", err)
    }
    defer f.Close()

    p := tea.NewProgram(model.InitialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }

}
