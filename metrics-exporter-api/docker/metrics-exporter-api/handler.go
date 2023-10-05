/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// Function to create route/handler for web request
func allHandlers() http.Handler {
	router := mux.NewRouter()
	// allowed CORS
	// more info https://en.wikipedia.org/wiki/Cross-origin_resource_sharing
	handler := cors.Default().Handler(router)

	router.HandleFunc("/", rootGet).Methods("GET")
	router.HandleFunc("/healthz", healthzGet).Methods("GET")
	// Openmetrics endpoints
	router.HandleFunc("/metrics", metricsGet).Methods("GET")
	router.HandleFunc("/metrics/device/{DeviceName}", deviceGet).Methods("GET")
	router.HandleFunc("/metrics/pci-addr/{PciAddr}", pciAddrGet).Methods("GET")
	// json metrics endpoints
	router.HandleFunc("/json/metrics", metricsGetJSON).Methods("GET")
	router.HandleFunc("/json/metrics/device/{DeviceName}", deviceGetJSON).Methods("GET")
	router.HandleFunc("/json/metrics/pci-addr/{PciAddr}", pciAddrGetJSON).Methods("GET")

	return handler
}

// this endpoint shows the uptime of the application
// it may be helpful to create probes in K8s
func healthzGet(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("uptime: %s\n", time.Since(time.Unix(0, StartupTime)))
	_, err := w.Write([]byte(response))

	if err != nil {
		log.Error(err)
	}
}

// Since nothing is to show on root handler
func rootGet(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("/ root uptime: %s\n", time.Since(time.Unix(0, StartupTime)))
	_, err := w.Write([]byte(response))

	if err != nil {
		log.Error(err)
	}
}
