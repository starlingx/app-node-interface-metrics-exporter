/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"

	"github.com/bsm/openmetrics"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

// endpoint to get all network metric of a node on which it is resides
// http://<hostname>:<port>/metrics
func metricsGet(w http.ResponseWriter, _ *http.Request) {

	openMetContent := allDevInfoOpenMet()
	w.Header().Set("Content-Type", OpenMetContentType)
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(openMetContent))

	if err != nil {
		log.Error(err)
	}

}

// endpoint to fetch metrics related to given network
// device by name
// http://<hostname>:<port>/device/<DeviceName>
func deviceGet(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	devName := params["DeviceName"]
	allDeviceStat := ListAllNetDev()

	devStats, ok := allDeviceStat[devName]

	// If the key exists
	if ok {
		devStatOpnFmt := devStatOpenMet(devStats)
		w.Header().Set("Content-Type", OpenMetContentType)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(devStatOpnFmt))

		if err != nil {
			log.Error(err)
		}
	} else {
		w.Header().Set("Content-Type", OpenMetContentType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", devName)

	}

}

// endpoint to fetch metrics related to given network
// device by pci addr
// http://<hostname>:<port>/pci-addr/<PciAddr>
func pciAddrGet(w http.ResponseWriter, r *http.Request) {

	found := false
	params := mux.Vars(r)
	PciAddr := params["PciAddr"]
	allDeviceStat := ListAllNetDev()

	for _, devStats := range allDeviceStat {

		if devStats.Pciaddr == PciAddr {
			found = true
			DevStatOpnFmt := devStatOpenMet(devStats)
			w.Header().Set("Content-Type", OpenMetContentType)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(DevStatOpnFmt))

			if err != nil {
				log.Error(err)
			}
		}
	}

	if !found {
		w.Header().Set("Content-Type", OpenMetContentType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Pci Addr Not found ", PciAddr)
	}

}

// function to register and update network device info
// and statistics to Openmetrics format
func allDevInfoOpenMet() string {

	// get all device info from lspci command
	allDeviceInfo := ListAllNetDev()
	// store stats data in openmetrics
	var openMetContent string

	for _, device := range allDeviceInfo {
		openMetContent += devStatOpenMet(device)
	}

	return openMetContent
}

func devStatOpenMet(devInfo DevInfo) string {
	var reg = openmetrics.NewRegistry()
	var openMetVar = regDevInfo(reg)

	openMetVar.With(
		devInfo.Name,
		devInfo.HardwareAddr,
		devInfo.Broadcast,
		// devInfo.duplex,
		devInfo.Alias,
		devInfo.OperState,
		devInfo.Pciaddr,
	)

	// convert all netlink.LinkStatistics Struct to map
	fields := reflect.TypeOf(*devInfo.Statistics)
	// valPtr := reflect.ValueOf(devInfo.Statistics)
	values := reflect.Indirect(reflect.ValueOf(devInfo.Statistics))
	// get the number of field for looping n times
	num := fields.NumField()
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		log.Debug("Type:", field.Type, ",", field.Name, "=", value, "\n")

		// register counter
		name := "network_interface_" + strcase.ToSnake(field.Name)
		var newInfo = reg.Counter(openmetrics.Desc{
			Name:   name,
			Help:   name,
			Labels: []string{"device"},
		})
		newInfo.With(
			devInfo.Name,
		).Add(float64(value.Uint()))
	}

	// create buffer to return
	var buf bytes.Buffer
	if _, err := reg.WriteTo(&buf); err != nil {
		panic(err)
	}
	// buffer.String()
	return buf.String()
}

// function to create Informational Openmetrics paramater / variable
// for Network Interface Device Info
func regDevInfo(reg *openmetrics.Registry) openmetrics.InfoFamily {

	var deviceInfo = reg.Info(openmetrics.Desc{
		Name:   "network_interface_device",
		Help:   "network_interface_device ",
		Labels: []string{"name", "address", "broadcast", "ifalias", "operstate", "pciaddr"},
	})

	return deviceInfo
}
