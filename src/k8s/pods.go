package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


func GetPods() [][]string {
    // fmt.Println("Connecting to cluster")

    // fmt.Println("Get Kubernetes pods")
    
    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
    // fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

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

    result := [][]string{}
    for _, pod := range pods.Items {
        var conditions string
        for _, component := range pod.Status.Conditions {
            conditions = conditions + string(component.Message) + ","
        }

        result = append(result, []string{
            pod.Name,
            pod.Namespace,
            fmt.Sprint(len(pod.Spec.Containers)),
            conditions,
        })
    }

    return result
}

