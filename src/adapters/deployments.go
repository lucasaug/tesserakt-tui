package adapters

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)

func GetDeploymentJson(
    clientset *kubernetes.Clientset,
    name, namespace string,
) string {
    deployment := k8s.GetDeployment(clientset, name, namespace)

    data, err := json.MarshalIndent(deployment, "", "    ")
    if err != nil {
        fmt.Printf("error getting json from deployment: %v\n", err)
        os.Exit(1)
    }

    return string(data)
}

func GetDeploymentRows(clientset *kubernetes.Clientset) []table.Row {
    rows := []table.Row {}
    for _, deployment := range k8s.GetDeployments(clientset) {
        rows = append(rows, table.Row{
            deployment.Name,
            deployment.Namespace,
        })
    }

    return rows
}


func GetDeploymentTable(clientset *kubernetes.Clientset) table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows([]table.Row{}),
    )
}
