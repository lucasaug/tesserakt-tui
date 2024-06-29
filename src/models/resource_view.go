package models

import (
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "k8s.io/client-go/kubernetes"

    "github.com/lucasaug/tesserakt-tui/src/commands"
    "github.com/lucasaug/tesserakt-tui/src/core"
)

type resourceView struct {
    clientset *kubernetes.Clientset

    width int
    height int

    style lipgloss.Style
    highlightedStyle lipgloss.Style
    resourceTables map[core.Resource]*table.Model
    contentViewport *viewport.Model

    resourceType core.Resource
    selectedResource map[core.Resource]*core.ResourceSelector
    itemIndex int
    highlighted bool
}

func InitialResourceViewModel() resourceView {
    viewStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder())
    highlightedStyle := lipgloss.NewStyle().
        BorderForeground(highlightColor).
        BorderStyle(lipgloss.NormalBorder())

    cv := viewport.New(0, 0)

    return resourceView{
        style: viewStyle,
        highlightedStyle: highlightedStyle,
        contentViewport: &cv,
        selectedResource: make(map[core.Resource]*core.ResourceSelector),
    }
}

func (r resourceView) Init() tea.Cmd {
    return nil
}

func (r resourceView) Update(msg tea.Msg) (resourceView, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case commands.K8sClientMsg:
        r.clientset = msg.Clientset
        cmd = commands.CreateLoadTables(r.clientset)

    case commands.TickMsg:
        cmd = commands.RefreshResourceList(r.clientset, &r.resourceType)

    case commands.ResourceChangeMsg:
        if (r.resourceType != msg.NewResource) {
            r.resourceType = msg.NewResource
            r.itemIndex = 0
            if (r.resourceTables[r.resourceType] != nil) {
                r.resourceTables[r.resourceType].SetCursor(0)
            }
            cmd = commands.RefreshResourceList(r.clientset, &r.resourceType)
        }

    case commands.LoadTablesMsg:
        r.resourceTables = msg.Tables

    case commands.RefreshResourceListMsg:
        if (len(r.resourceTables) != 0) {
            r.resourceTables[msg.Resource].SetRows(msg.Rows)
        }

    case commands.ResourceDetailsMsg:
        r.selectedResource[r.resourceType] = &msg.Value

    case tea.KeyMsg:
        switch msg.String() {

        case "enter":
            name := r.resourceTables[r.resourceType].
                Rows()[r.itemIndex][0]
            namespace := r.resourceTables[r.resourceType].
                Rows()[r.itemIndex][1]

            return r, commands.ResourceDetails(
                *r.clientset,
                r.resourceType,
                name,
                namespace,
            )

        case "k", "up":
            if (r.itemIndex > 0 &&
                r.selectedResource[r.resourceType] == nil) {
                r.itemIndex--
            }

        case "j", "down":
            nrows := len(r.resourceTables[r.resourceType].Rows())
            if (r.itemIndex < nrows - 1 &&
                r.selectedResource[r.resourceType] == nil) {
                r.itemIndex++
            }

        case "esc":
            r.selectedResource[r.resourceType] = nil

        }
    }

    if (r.selectedResource[r.resourceType] != nil) {
        var cvCmd tea.Cmd
        *r.contentViewport, cvCmd = r.contentViewport.Update(msg)
        cmd = tea.Batch(cvCmd, cmd)
    } else if (len(r.resourceTables) != 0) {
        var tableCmd tea.Cmd
        *r.resourceTables[r.resourceType], tableCmd =
            r.resourceTables[r.resourceType].Update(msg)
        cmd = tea.Batch(tableCmd, cmd)
    }

    return r, cmd
}

func (r resourceView) View() string {
    if (r.selectedResource[r.resourceType] != nil) {
        r.contentViewport.Width = r.width
        r.contentViewport.Height = r.height
        r.contentViewport.SetContent(r.selectedResource[r.resourceType].Data)

        if (r.highlighted) {
            return r.highlightedStyle.Render(
                r.contentViewport.View(),
            )
        }
        return r.style.Render(r.contentViewport.View())
    }

    if (len(r.resourceTables) == 0) { return "" }

    r.resourceTables[r.resourceType].SetWidth(r.width)
    r.resourceTables[r.resourceType].SetHeight(r.height)

    if (r.highlighted) {
        return r.highlightedStyle.Render(
            r.resourceTables[r.resourceType].View(),
        )
    }

    return r.style.Render(r.resourceTables[r.resourceType].View())
}

func (r resourceView) Focus() {
    r.resourceTables[r.resourceType].Focus()
}

func (r resourceView) Blur() {
    r.resourceTables[r.resourceType].Blur()
}

func (r *resourceView) SetSize(width int, height int) {
    r.width = width
    r.height = height
}

func (r *resourceView) SetHighlight(highlighted bool) {
    r.highlighted = highlighted
}
