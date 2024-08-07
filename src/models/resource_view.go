package models

import (
    "bytes"

    "github.com/alecthomas/chroma/quick"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "k8s.io/client-go/kubernetes"

    "github.com/lucasaug/tesserakt-tui/src/commands"
    "github.com/lucasaug/tesserakt-tui/src/controllers"
    "github.com/lucasaug/tesserakt-tui/src/k8s"
)

const DEFAULT_RESOURCE = k8s.Pod

type resourceView struct {
    clientset *kubernetes.Clientset

    width int
    height int

    style lipgloss.Style
    highlightedStyle lipgloss.Style
    currentTable *table.Model
    contentViewport *viewport.Model

    resourceType k8s.ResourceType
    selectedResource map[k8s.ResourceType]*k8s.ResourceSelector
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
        selectedResource: make(map[k8s.ResourceType]*k8s.ResourceSelector),
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
        r.resourceType = DEFAULT_RESOURCE
        r.itemIndex = 0

        currentTable := table.New(
            table.WithColumns(controllers.GetColumns(r.resourceType)),
            table.WithRows([]table.Row{}),
        )
        r.currentTable = &currentTable
        r.currentTable.SetCursor(0)

        cmd = commands.RefreshResourceList(r.clientset, &r.resourceType)

    case commands.TickMsg:
        cmd = commands.RefreshResourceList(r.clientset, &r.resourceType)

    case commands.ResourceChangeMsg:
        if (r.resourceType != msg.NewResource) {
            r.resourceType = msg.NewResource
            r.itemIndex = 0
            currentTable := table.New(
                table.WithColumns(controllers.GetColumns(r.resourceType)),
                table.WithRows([]table.Row{}),
            )
            r.currentTable = &currentTable
            cmd = commands.RefreshResourceList(r.clientset, &r.resourceType)
        }

    case commands.RefreshResourceListMsg:
        if (r.currentTable != nil) {
            r.currentTable.SetRows(msg.Rows)
        }

    case commands.ResourceDetailsMsg:
        r.selectedResource[r.resourceType] = &msg.Value

    case tea.KeyMsg:
        switch msg.String() {

        case "enter":
            if (r.selectedResource[r.resourceType] == nil) {
                name := r.currentTable.Rows()[r.itemIndex][0]
                namespace := r.currentTable.Rows()[r.itemIndex][1]

                return r, commands.ResourceDetails(
                    *r.clientset,
                    r.resourceType,
                    name,
                    namespace,
                )
            }

        case "k", "up":
            if (r.itemIndex > 0 &&
                r.selectedResource[r.resourceType] == nil) {
                r.itemIndex--
            }

        case "j", "down":
            nrows := len(r.currentTable.Rows())
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
    } else if (r.currentTable != nil) {
        var tableCmd tea.Cmd
        *r.currentTable, tableCmd = r.currentTable.Update(msg)
        cmd = tea.Batch(tableCmd, cmd)
    }

    return r, cmd
}

func (r resourceView) View() string {
    if (r.selectedResource[r.resourceType] != nil) {
        r.contentViewport.Width = r.width
        r.contentViewport.Height = r.height
        data := r.selectedResource[r.resourceType].Data
        var b bytes.Buffer
        quick.Highlight(&b, data, "json", "terminal", "catpuccin-mocha")
        highlightedYAML := b.String()
        r.contentViewport.SetContent(highlightedYAML)

        if (r.highlighted) {
            return r.highlightedStyle.Render(
                r.contentViewport.View(),
            )
        }
        return r.style.Render(r.contentViewport.View())
    }

    if (r.currentTable == nil) { return "" }

    r.currentTable.SetWidth(r.width)
    r.currentTable.SetHeight(r.height)

    if (r.highlighted) {
        return r.highlightedStyle.Render(
            r.currentTable.View(),
        )
    }

    return r.style.Render(r.currentTable.View())
}

func (r resourceView) Focus() {
    r.currentTable.Focus()
}

func (r resourceView) Blur() {
    r.currentTable.Blur()
}

func (r *resourceView) SetSize(width int, height int) {
    r.width = width
    r.height = height
}

func (r *resourceView) SetHighlight(highlighted bool) {
    r.highlighted = highlighted
}
