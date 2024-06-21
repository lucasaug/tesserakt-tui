package adapters

import (
	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)


func GetDeploymentTable(clientset *kubernetes.Clientset) table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    }

    rows := []table.Row {}
    for _, deployment := range k8s.GetDeployments(clientset) {
        rows = append(rows, table.Row{
            deployment.Name,
            deployment.Namespace,
        })
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
    )
}
