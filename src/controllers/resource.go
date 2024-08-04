package controllers

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	"github.com/lucasaug/tesserakt-tui/src/k8s"
)

var resourceHandlers = map[k8s.ResourceType]k8s.ResourceHandler {
    k8s.Pod: k8s.PodHandler{},
    k8s.ReplicaSet: k8s.ReplicaSetHandler{},
    k8s.StatefulSet: k8s.StatefulSetHandler{},
    k8s.DaemonSet: k8s.DaemonSetHandler{},
    k8s.Job: k8s.JobHandler{},
    k8s.CronJob: k8s.CronJobHandler{},
    k8s.Deployment: k8s.DeploymentHandler{},
    k8s.Ingress: k8s.IngressHandler{},
}

func GetColumns(resourceType k8s.ResourceType) []table.Column  {
    return resourceHandlers[resourceType].Columns()
}

func Get(
    clientset *kubernetes.Clientset,
    resourceType k8s.ResourceType,
    name, namespace string,
) (string, error) {

    handler:= resourceHandlers[resourceType]

    resource := k8s.GetResource(
        clientset,
        handler,
        name,
        namespace,
    )

    data, err := json.MarshalIndent(resource, "", "    ")
    if err != nil {
        return "", err
    }

    return string(data), nil
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
        table.WithColumns(GetColumns(resourceType)),
        table.WithRows([]table.Row{}),
    )
}
