package adapters

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)


func GetPodTable() table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 10 },
        { Title: "Namespace", Width: 10 },
        { Title: "Num of containers", Width: 20 },
        { Title: "Conditions", Width: 20 },
    }

    rows := []table.Row {}
    for _, pod := range k8s.GetPods() {
        var conditions string
        for _, component := range pod.Status.Conditions {
            conditions = conditions + string(component.Message) + ","
        }

        rows = append(rows, table.Row{
            pod.Name,
            pod.Namespace,
            fmt.Sprint(len(pod.Spec.Containers)),
            conditions,
        })
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
    )
}
