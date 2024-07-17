package controllers

import (
    "github.com/charmbracelet/bubbles/table"
    "k8s.io/client-go/kubernetes"

    "github.com/lucasaug/tesserakt-tui/src/k8s"
)

var resourceHandlers = map[k8s.ResourceType]k8s.ResourceHandler {
    k8s.Pod: k8s.PodHandler{},
    k8s.Deployment: k8s.DeploymentHandler{},
    k8s.Ingress: k8s.IngressHandler{},
}

var ResourceToColumns = map[k8s.ResourceType][]table.Column {
    k8s.Pod: {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
        { Title: "Container count", Width: 15 },
        { Title: "Phase", Width: 20 },
    },
    k8s.Deployment: {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    },
    k8s.Ingress: {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    },
}

func GetRows(
    clientset *kubernetes.Clientset,
    resourceType k8s.ResourceType,
) ([]table.Row, error) {
    rows := []table.Row{}
    items, err := resourceHandlers[resourceType].List(clientset, "")

    if err != nil {
        return rows, err
    }

    for _, item := range items {
        row := table.Row{}
        row = append(row, item.Values()...)

        rows = append(rows, row)
    }

    return rows, nil
}

func GetTable(resourceType k8s.ResourceType) table.Model {
    return table.New(
        table.WithColumns(ResourceToColumns[resourceType]),
        table.WithRows([]table.Row{}),
    )
}
