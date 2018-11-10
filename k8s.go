package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var serializer = json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme,
	scheme.Scheme)

// Saves serialized object in a file
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

	fmt.Println("\tStored", name, "in", filePath)
	return nil
}

// Fetch PVs and save in the file
func storePersistentVolumes(client corev1.CoreV1Interface, path string, serializer *json.Serializer) {
	fmt.Println("Storing PersistentVolumes")

	// Make directory for persistent volumes
	pathPrefix := filepath.Join(path, "persistentvolumes")
	os.Mkdir(pathPrefix, os.ModePerm)
	fmt.Println("\tCreated", pathPrefix)

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
func storeGlobals(client corev1.CoreV1Interface, path string) {
	// Make directory for globals
	globalsPath := filepath.Join(path, "globals")
	os.Mkdir(globalsPath, os.ModePerm)
	fmt.Println("Created", globalsPath)

	storePersistentVolumes(client, globalsPath, serializer)
	// storeSecurityContextConstraints(client, globalsPath, serializer)
	// storeClusterRoles(client, globalsPath, serializer)
	// storeClusterRoleBindings(client, globalsPath, serializer)
}

// Get a list of namespaces
func getNamespaces(client corev1.CoreV1Interface) []string {
	fmt.Println("Fetching namespaces")
	nsList, err := client.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	result := make([]string, len(nsList.Items))
	for i, ns := range nsList.Items {
		result[i] = ns.Name
	}
	return result
}

// Store namespace objects
func storeNamespace(client corev1.CoreV1Interface, namespace string, path string) {
	ns, err := client.Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		// TODO collect errs here instead of panicking
		panic(err)
	}

	name := ns.Name

	// TODO: serialize namespace
	/* var buf bytes.Buffer
	err = serializer.Encode(&ns, &buf)
	if err != nil {
		// TODO collect errs here instead of panicking
		panic(err)
	}
	err = storeObject(&buf, name, path)
	if err != nil {
		// TODO collect errs here instead of panicking
		panic(err)
	} */

	namespacePath := filepath.Join(path, name)
	os.Mkdir(namespacePath, os.ModePerm)
	fmt.Println("Created", namespacePath)

	// TODO: secrets are forbidden
	// storeSecrets(client, namespacePath, serializer, name)
}

// Store secrets
func storeSecrets(client corev1.CoreV1Interface, pathPrefix string, serializer *json.Serializer, namespace string) {
	fmt.Println("\tStoring secrets")
	secretsList, err := client.Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		// TODO collect errs here instead of panicking
		panic(err)
	}

	// Save each volume
	for _, secret := range secretsList.Items {
		name := secret.Name
		var buf bytes.Buffer
		err = serializer.Encode(&secret, &buf)
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
