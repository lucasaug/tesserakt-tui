package adapters

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"k8s.io/client-go/kubernetes"

	k8s "github.com/lucasaug/tesserakt-tui/src/k8s"
)

func GetIngressJson(
    clientset *kubernetes.Clientset,
    name, namespace string,
) string {
    ingress := k8s.GetIngress(clientset, name, namespace)

    data, err := json.MarshalIndent(ingress, "", "    ")
    if err != nil {
        fmt.Printf("error getting json from ingress: %v\n", err)
        os.Exit(1)
    }

    return string(data)
}


func GetIngressRows(clientset *kubernetes.Clientset) []table.Row {
    rows := []table.Row {}
    for _, ingress := range k8s.GetIngresses(clientset) {
        rows = append(rows, table.Row{
            ingress.Name,
            ingress.Namespace,
        })
    }

    return rows
}

func GetIngressTable(clientset *kubernetes.Clientset) table.Model {
    columns := []table.Column {
        { Title: "Name", Width: 40 },
        { Title: "Namespace", Width: 20 },
    }

    return table.New(
        table.WithColumns(columns),
        table.WithRows([]table.Row{}),
    )
}
