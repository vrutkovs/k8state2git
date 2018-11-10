package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

var serializer = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme,
	scheme.Scheme)

// Clone the repo and return the path
func cloneRepo() string {
	return "/repo"
}

func storeObject(buffer *bytes.Buffer, name string, pathPrefix string) error {
	filePath, err := filepath.Abs(filepath.Join(pathPrefix, fmt.Sprintf("%s.yml", name)))
	if err != nil {
		return err
	}

	// Create file
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = buffer.WriteTo(w)
	if err != nil {
		return err
	}
	w.Flush()

	fmt.Println("\tStored ", name, " in ", filePath)
	return nil
}

func storePersistentVolumes(client corev1.CoreV1Interface, path string, serializer *json.Serializer) {
	fmt.Println("Storing PersistentVolumes")

	// Make directory for persistent volumes
	pathPrefix := filepath.Join(path, "persistentvolumes")
	os.Mkdir(pathPrefix, os.ModePerm)
	fmt.Println("\tCreated ", pathPrefix)

	// List volumes
	pvList, err := client.PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Save each volume
	for _, pv := range pvList.Items {
		name := pv.Name
		var buf bytes.Buffer
		err = serializer.Encode(&pv, &buf)
		if err != nil {
			// TODO collect errs here instead of panicking
			panic(err)
		}
		err = storeObject(&buf, name, pathPrefix)
		if err != nil {
			// TODO collect errs here instead of panicking
			panic(err)
		}
	}
}

// Store global k8s objects
func storeGlobals(client corev1.CoreV1Interface, path string, serializer *json.Serializer) {
	// Make directory for globals
	globalsPath := filepath.Join(path, "globals")
	os.Mkdir(globalsPath, os.ModePerm)
	fmt.Println("Created ", globalsPath)

	storePersistentVolumes(client, globalsPath, serializer)

}

// Get a list of namespaces
func getNamespaces(client corev1.CoreV1Interface) []string {
	return make([]string, 3)
}

// Store namespace objects
func storeNamespaces(namespace string, path string) {

}

// Make a git commit
func gitCommit(path string) {

}

// Push git commit
func gitPush(path string) {

}

func saveClusterState(clientset *kubernetes.Clientset, serializer *json.Serializer) {
	var client = clientset.CoreV1()
	path := cloneRepo()
	storeGlobals(client, path, serializer)
	namespaces := getNamespaces(client)
	for _, namespace := range namespaces {
		storeNamespaces(namespace, path)
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

	saveClusterState(clientset, serializer)

}
