package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Connect to k8s and save all objects
func saveClusterState(clientset *kubernetes.Clientset) {
	var client = clientset.CoreV1()
	path, err := cloneRepo()
	if err != nil {
		panic(err.Error())
	}
	if err := gitConfig(path); err != nil {
		panic(err.Error())
	}
	if err := cleanRepo(path); err != nil {
		panic(err.Error())
	}

	storeGlobals(client, path)
	namespaces := getNamespaces(client)
	for _, namespace := range namespaces {
		storeNamespace(client, namespace, path)
	}

	gitCommit(path)
	gitPush(path)

}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// connect to k8s
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// Read k8k8state2git config here

	saveClusterState(clientset)

}
