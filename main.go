package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"

    tea "github.com/charmbracelet/bubbletea"
)

type Resource string

const (
    Pod        Resource = "Pod"
    Deployment Resource = "Deployment"
    Ingress    Resource = "Ingress"
)

type model struct {
    currentResource Resource
    names           []string
    currentIndex    int
}

func initialModel() model {
    return model{
        currentResource: "bruh",
        names: getPods(),
        currentIndex: 0,
    }
}

func (m model) Init() tea.Cmd {
    return nil;
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:

        switch msg.String() {

        case "q":
            return m, tea.Quit

        case "up", "k":
            if m.currentIndex < len(m.names) {
                m.currentIndex++
            }

        case "down", "j":
            if m.currentIndex > 0 {
                m.currentIndex--
            }

        }
    }

    return m, nil
}

func (m model) View() string {
    s := "List"

    for i, name := range m.names {
        cursor := " "
        if m.currentIndex == i {
            cursor = ">"
        }

        s += fmt.Sprintf("%s [%s]\n", cursor, name)
    }

    return s
}


func getPods() []string {
    fmt.Println("Get Kubernetes pods")

    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
    fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

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

    result := []string{}
    for _, pod := range pods.Items {
        result = append(result, pod.Name)
    }

    return result
}

func main() {
    fmt.Println("Connecting to cluster")

    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Fuck")
        os.Exit(1)
    }

}
