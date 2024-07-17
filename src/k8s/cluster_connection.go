package k8s

import (
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func GetClientSet() *kubernetes.Clientset {
    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

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

    return clientset
}

