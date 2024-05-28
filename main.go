package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/charmbracelet/bubbles/table"
	// "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Resource string

const (
    Pod        Resource = "Pod"
    Deployment Resource = "Deployment"
    Ingress    Resource = "Ingress"
)

type item struct {
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.title }
func (i item) FilterValue() string { return i.title }

type model struct {
    currentResource Resource
    currentIndex    int
    resourceList    table.Model
    table           table.Model
    width           int
    height          int
    tableStyle      lipgloss.Style
    listStyle       lipgloss.Style
}

func initialModel() model {
    resourceItems := []table.Row {
        []string{ string(Pod) },
        []string{ string(Deployment) },
        []string{ string(Ingress) },
    }

    columns := []table.Column {
        { Title: "Name", Width: 10 },
        { Title: "Namespace", Width: 10 },
        { Title: "Num of containers", Width: 20 },
    }

    rows := []table.Row {}
    for _, pod := range getPods() {
        rows = append(rows, pod)
    }

    borderColor := lipgloss.Color("36")
    tableStyle := lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
    listStyle := lipgloss.NewStyle().BorderForeground(borderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(20)

    return model{
        currentResource: "Pod",
        currentIndex: 0,
        table: table.New(
            table.WithColumns(columns),
            table.WithRows(rows),
            table.WithFocused(true),
            table.WithHeight(7),
        ),
        resourceList: table.New(
            table.WithColumns([]table.Column {{ Title: "Resources", Width: 10 }}),
            table.WithRows(resourceItems),
            table.WithFocused(true),
            table.WithHeight(7),
        ),
        tableStyle: tableStyle,
        listStyle: listStyle,
    }
}

func (m model) Init() tea.Cmd {
    return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case tea.KeyMsg:

        switch msg.String() {

        case "ctrl+c", "q":
            return m, tea.Quit

        // case "up", "k":
        //     if m.currentIndex < len(m.names) {
        //         m.currentIndex++
        //     }
        //
        // case "down", "j":
        //     if m.currentIndex > 0 {
        //         m.currentIndex--
        //     }
        //
        }
    }

    return m, nil
}

func (m model) View() string {
    return lipgloss.Place(
        m.width,
        m.height,
        lipgloss.Center,
        lipgloss.Center,
        lipgloss.JoinHorizontal(
            lipgloss.Top,
            m.listStyle.Render(m.resourceList.View()),
            m.tableStyle.Render(m.table.View()),
        ),
    )
}


func getPods() [][]string {
    // fmt.Println("Connecting to cluster")

    // fmt.Println("Get Kubernetes pods")

    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
    // fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

    kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        fmt.Printf("error getting Kubernetes config: %v\n", err)
        os.Exit(1)
    }

    clientset, err := kubernetes.NewForConfig(kubeConfig)
    if err != nil {
        fmt.Printf("error getting Kubernetes clientset: %v\n", err)
        os.Exit(1)
    }

    pods, err := clientset.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{})
    if err != nil {
        fmt.Printf("error getting pods: %v\n", err)
        os.Exit(1)
    }

    result := [][]string{}
    for _, pod := range pods.Items {
        result = append(result, []string{
            pod.Name,
            pod.Namespace,
            fmt.Sprint(len(pod.Spec.Containers)),
        })
    }

    return result
}

func main() {
    f, err := tea.LogToFile("debug.log", "debug")
    if err != nil {
        fmt.Println("Error creating log file:", err)
    }
    defer f.Close()

    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error running program:", err)
        os.Exit(1)
    }

}
