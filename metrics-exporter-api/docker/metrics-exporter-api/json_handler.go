/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Response Struct
type Response struct {
	Devices PfDevices
	VfPod   VfpodInfo
	VfDev   vfDevices
}

// endpoint to get network metric of a node on which it is resides
// http://<hostname>:<port>/json/metrics
func (m *metricsHandler) metricsGetJSON(w http.ResponseWriter, _ *http.Request) {
	defer m.recoverFromPanic(w)
	res := Response{
		Devices: m.ListAllNetDev(),
		VfPod:   m.fetchVfPodInfo(),
	}
	// convert the map to a JSON encoded byte slice
	handleJsonResponse(res, w)

}

// endpoint to fetch metrics related to given network
// device by name
// http://<hostname>:<port>/json/metrics/device/<DeviceName>
func (m *metricsHandler) deviceGetJSON(w http.ResponseWriter, r *http.Request) {
	defer m.recoverFromPanic(w)
	params := mux.Vars(r)
	DeviceName := params["DeviceName"]
	res := m.deviceByProp("Name", DeviceName)

	if res.Devices == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", DeviceName)
		return
	}

	handleJsonResponse(res, w)

}

// endpoint to fetch metrics related to given network
// device by pci addr
// http://<hostname>:<port>/json/metrics/pci-addr/<PciAddr>
func (m *metricsHandler) pciAddrGetJSON(w http.ResponseWriter, r *http.Request) {
	defer m.recoverFromPanic(w)
	params := mux.Vars(r)
	PciAddr := params["PciAddr"]

	// first check into Physical devices
	res := m.deviceByProp("Pciaddr", PciAddr)

	// if not found in physical devices search in vf

	if res.Devices == nil && res.VfDev == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", PciAddr)
		return
	}

	handleJsonResponse(res, w)

}

// Showing Up time and version on root handler
func rootGet(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"uptime":  time.Since(time.Unix(0, StartupTime)).String(),
		"version": Version,
		"build":   Build,
	}

	handleJsonResponse(response, w)
}

func handleJsonResponse(res interface{}, w http.ResponseWriter) {
	// In case of no errors, status code defaulted to 200
	w.Header().Set("Content-Type", "application/json")

	// convert the map to a JSON encoded byte slice
	jsonContent, mErr := json.Marshal(res)
	if mErr != nil {
		log.Error(mErr)
		http.Error(w, fmt.Sprintf("Error encoding json res: %s", mErr), http.StatusInternalServerError)
		return
	}

	// Write method sets the statusCode to 200 if we don't encounter any error
	_, err := w.Write(jsonContent)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("Error writing response body: %s", err), http.StatusInternalServerError)
		return
	}
}
