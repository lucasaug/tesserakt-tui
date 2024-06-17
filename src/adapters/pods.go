package adapters

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)


func GetPodTable() table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 50 },
        { Title: "Namespace", Width: 20 },
        { Title: "Container count", Width: 15 },
        { Title: "Phase", Width: 20 },
    }

    rows := []table.Row {}
    for _, pod := range k8s.GetPods() {
        rows = append(rows, table.Row{
            pod.Name,
            pod.Namespace,
            fmt.Sprint(len(pod.Spec.Containers)),
            string(pod.Status.Phase),
        })
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
    )
}
