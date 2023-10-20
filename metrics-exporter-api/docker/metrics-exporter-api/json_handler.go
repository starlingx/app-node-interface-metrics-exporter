/*
 Copyright (c) 2023 Wind River Systems, Inc.

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
func metricsGetJSON(w http.ResponseWriter, _ *http.Request) {

	res := Response{
		Devices: ListAllNetDev(),
		VfPod:   fetchVfPodInfo(),
	}
	// convert the map to a JSON encoded byte slice
	jsonContent, mErr := json.Marshal(res)
	if mErr != nil {
		log.Error("Error: ", mErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(jsonContent)

	if err != nil {
		log.Error(err)
	}

}

// endpoint to fetch metrics related to given network
// device by name
// http://<hostname>:<port>/json/metrics/device/<DeviceName>
func deviceGetJSON(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	DeviceName := params["DeviceName"]
	res := deviceByProp("Name", DeviceName)

	if res.Devices == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", DeviceName)
		return
	}

	// convert the map to a JSON encoded byte slice
	jsonContent, mErr := json.Marshal(res)
	if mErr != nil {
		log.Error(mErr)
		return
	}

	// convert the byte slice to a string
	// jsonString := string(jsonContent)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(jsonContent)

	if err != nil {
		log.Error(err)
	}

}

// endpoint to fetch metrics related to given network
// device by pci addr
// http://<hostname>:<port>/json/metrics/pci-addr/<PciAddr>
func pciAddrGetJSON(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	PciAddr := params["PciAddr"]

	// first check into Physical devices
	res := deviceByProp("Pciaddr", PciAddr)

	// if not found in physical devices search in vf

	if res.Devices == nil && res.VfDev == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", PciAddr)
		return
	}

	// convert the map to a JSON encoded byte slice
	jsonContent, mErr := json.Marshal(res)
	if mErr != nil {
		log.Error(mErr)
		return
	}

	// convert the byte slice to a string
	// jsonString := string(jsonContent)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write(jsonContent)

	if err != nil {
		log.Error(err)
	}

}

// Showing Up time and version on root handler
func rootGet(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"uptime":  time.Since(time.Unix(0, StartupTime)).String(),
		"version": Version,
		"build":   Build,
	}

	// convert the map to a JSON encoded byte slice
	jsonContent, mErr := json.Marshal(response)
	if mErr != nil {
		log.Error(mErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(jsonContent)

	if err != nil {
		log.Error(err)
	}
}
