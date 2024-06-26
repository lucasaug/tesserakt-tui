package adapters

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)

func GetPodJson(
    clientset *kubernetes.Clientset,
    name, namespace string,
) string {
    pod := k8s.GetPod(clientset, name, namespace)

    data, err := json.MarshalIndent(pod, "", "    ")
    if err != nil {
        fmt.Printf("error getting json from pod: %v\n", err)
        os.Exit(1)
    }

    return string(data)
}

func GetPodRows(clientset *kubernetes.Clientset) []table.Row {
    rows := []table.Row {}
    for _, pod := range k8s.GetPods(clientset) {
        rows = append(rows, table.Row{
            pod.Name,
            pod.Namespace,
            fmt.Sprint(len(pod.Spec.Containers)),
            string(pod.Status.Phase),
        })
    }

    return rows
}

func GetPodTable(clientset *kubernetes.Clientset) table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 15 },
        { Title: "Container count", Width: 15 },
        { Title: "Phase", Width: 20 },
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows([]table.Row{}),
    )
}
