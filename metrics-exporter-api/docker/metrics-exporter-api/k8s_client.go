/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// VfpodInfo structure
type VfpodInfo map[string]map[string]string

func getClientSet() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		log.Error(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err.Error())
		// here panic is required since it not able to communicate woth k8s api
		panic(err.Error())
	}

	return clientset

}

// Class to contain all methods which deal with VF kube pods
type VfPodClass interface {
	fetchVfPodInfo() VfpodInfo
}

// This struct encapsulates the dependencies of fetchVfPodInfo
// aids in creating mocks for these dependencies and test functions
type VfPodReceiver struct {
	// empty
}

// get all the pods which are using virtual function
func (v *VfPodReceiver) fetchVfPodInfo() VfpodInfo {
	clientset := getClientSet()
	// init map which will contain pod data
	vfPods := VfpodInfo{}

	// from Downward API we can get node name
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + os.Getenv("MY_NODE_NAME"),
	})

	if err != nil {
		log.Error(err.Error())
	}
	// loop all pods to check vf
	for _, pod := range pods.Items {
		// loop all container in pod to check vf
		for _, container := range pod.Spec.Containers {
			for key := range container.Resources.Requests {
				if strings.Contains(key.String(), vfIdentifier) {
					// Getting some information from pod annotations
					// from k8s.v1.cni.cncf.io/networks-status:
					for AnnotationName, AnnotationVal := range pod.GetAnnotations() {
						if strings.Contains(AnnotationName, "k8s.v1.cni.cncf.io/network-status") {
							annoMap := castStr2Map(AnnotationVal)
							// since we got annotation map we need to find all the network with interface
							for index, netInterface := range annoMap {
								// check if interface exists
								_, exists := netInterface["interface"]
								if exists {
									pciDet := getValfromInterface(netInterface["device-info"], "pci")
									pciAddr := getValfromInterface(pciDet, "pci-address")

									vfPods[fmt.Sprintf("%s-%d", container.Name, index)] = map[string]string{
										"Namespace":    pod.GetNamespace(),
										"Pod":          pod.Name,
										"Container":    container.Name,
										"Pciaddr":      fmt.Sprintf("%s", pciAddr),
										"HardwareAddr": fmt.Sprintf("%s", netInterface["mac"]),
										"Resource":     key.String(),
										"Device":       fmt.Sprintf("%s", netInterface["interface"]),
										"Vf":           fmt.Sprintf("%d", index),
									}
								}
							}
						}
					}
				}

			}

		}
	}
	return vfPods
}
