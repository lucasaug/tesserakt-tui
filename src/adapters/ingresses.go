package adapters

import (
	"github.com/charmbracelet/bubbles/table"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)


func GetIngressTable() table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 50 },
        { Title: "Namespace", Width: 20 },
    }

    rows := []table.Row {}
    for _, ingress := range k8s.GetIngresses() {
        rows = append(rows, table.Row{
            ingress.Name,
            ingress.Namespace,
        })
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
    )
}
