package models

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/adapters"
)

type resourceSelector struct {
    name string
    namespace string
    resourceType Resource
    editable bool
    data string
}

type resourceView struct {
    clientset *kubernetes.Clientset
    itemIndex int

    style lipgloss.Style
    highlightedStyle lipgloss.Style

    width int
    height int

    resourceType Resource
    resourceTables map[Resource]*table.Model
    contentViewport *viewport.Model

    selectedResource *resourceSelector
    highlighted bool
}

func InitialResourceViewModel() resourceView {
    tableStyle := lipgloss.NewStyle().
        BorderForeground(borderColor).
        BorderStyle(lipgloss.NormalBorder())
    highlightedStyle := lipgloss.NewStyle().
        BorderForeground(highlightColor).
        BorderStyle(lipgloss.NormalBorder())

    return resourceView{
        style: tableStyle,
        highlightedStyle: highlightedStyle,
    }
}

func (r *resourceView) SetResource(res Resource) {
    r.resourceType = res
}

const TICK_INTERVAL = time.Millisecond * 500

type TickMsg time.Time

func tick() tea.Cmd {
    return tea.Tick(TICK_INTERVAL, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

type RefreshRowsMsg struct {
    rows []table.Row
    resource Resource
}
type EmptyMsg struct {}

func (r resourceView) refreshList() tea.Msg {
    if (r.clientset == nil) { return EmptyMsg{} }

    if r.resourceType == Pod {
        return RefreshRowsMsg{
            rows: adapters.GetPodRows(r.clientset),
            resource: r.resourceType,
        }
    } else if r.resourceType == Deployment {
        return RefreshRowsMsg{
            rows: adapters.GetDeploymentRows(r.clientset),
            resource: r.resourceType,
        }
    } else if r.resourceType == Ingress {
        return RefreshRowsMsg{
            rows: adapters.GetIngressRows(r.clientset),
            resource: r.resourceType,
        }
    }

    return RefreshRowsMsg{
        rows: []table.Row{},
        resource: r.resourceType,
    }
}

type ResourceDetailsMsg struct { value resourceSelector }

func createResourceDetails(
    clientset kubernetes.Clientset,
    resourceType Resource,
    name, namespace string,
) tea.Cmd {
    return func() tea.Msg {
        var data string
        if (resourceType == Pod) {
            // data = adapters.GetPodJson(&clientset, name, namespace)
        } else if (resourceType == Deployment) {
            data = adapters.GetDeploymentJson(&clientset, name, namespace)
        } else if (resourceType == Ingress) {
            // data = adapters.GetIngressJson(&clientset, name, namespace)
        }

        return ResourceDetailsMsg{
            value: resourceSelector{
                name: name,
                namespace: namespace,
                resourceType: resourceType,
                editable: false,
                data: data,
            },
        }
    }
}

type LoadTablesMsg struct { tables map[Resource]*table.Model }
func createLoadTables(clientset *kubernetes.Clientset) tea.Cmd {
    return func() tea.Msg {
        podTable := adapters.GetPodTable(clientset)
        deploymentTable := adapters.GetDeploymentTable(clientset)
        ingressTable := adapters.GetIngressTable(clientset)

        tables := map[Resource]*table.Model {
            Pod: &podTable,
            Deployment: &deploymentTable,
            Ingress: &ingressTable,
        }

        return LoadTablesMsg{tables: tables}
    }
}

func (r resourceView) Init() tea.Cmd {
    return tick()
}

func (r resourceView) Update(msg tea.Msg) (resourceView, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case k8sClientMsg:
        r.clientset = msg.clientset
        return r, createLoadTables(r.clientset)

    case TickMsg:
        return r, tea.Batch(r.refreshList, tick())

    case LoadTablesMsg:
        r.resourceTables = msg.tables

    case RefreshRowsMsg:
        if (len(r.resourceTables) != 0) {
            r.resourceTables[msg.resource].SetRows(msg.rows)
        }

    case ResourceDetailsMsg:
        r.selectedResource = &msg.value

    case tea.KeyMsg:
        switch msg.String() {

        case "enter":
            name := r.resourceTables[r.resourceType].Rows()[r.itemIndex][0]
            namespace := r.resourceTables[r.resourceType].Rows()[r.itemIndex][1]
            return r, createResourceDetails(*r.clientset, r.resourceType, name, namespace)

        case "k", "up":
            if (r.itemIndex > 0) {
                r.itemIndex--
            }

        case "j", "down":
            if (r.itemIndex < len(r.resourceTables[r.resourceType].Rows()) - 1) {
                r.itemIndex++
            }

        }
    }

    if (len(r.resourceTables) != 0) {
        *r.resourceTables[r.resourceType], cmd = r.resourceTables[r.resourceType].Update(msg)
    }

    return r, cmd
}

func (r resourceView) View() string {
    if (r.selectedResource != nil) {
        r.contentViewport.SetContent(r.selectedResource.data)
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
