/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

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
	metricsHandler := newMetricsHandler()
	router.HandleFunc("/", rootGet).Methods("GET")
	router.HandleFunc("/healthz", healthzGet).Methods("GET")
	// Openmetrics endpoints
	router.HandleFunc("/metrics", metricsHandler.metricsGet).Methods("GET")
	// router.Methods("GET").Path("/metrics").H
	router.HandleFunc("/metrics/device/{DeviceName}", metricsHandler.deviceGet).Methods("GET")
	router.HandleFunc("/metrics/pci-addr/{PciAddr}", metricsHandler.pciAddrGet).Methods("GET")
	// json metrics endpoints
	router.HandleFunc("/json/metrics", metricsHandler.metricsGetJSON).Methods("GET")
	router.HandleFunc("/json/metrics/device/{DeviceName}", metricsHandler.deviceGetJSON).Methods("GET")
	router.HandleFunc("/json/metrics/pci-addr/{PciAddr}", metricsHandler.pciAddrGetJSON).Methods("GET")

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
