/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/bsm/openmetrics"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

// endpoint to get all network metric of a node on which it is resides
// http://<hostname>:<port>/metrics
func (m *metricsHandler) metricsGet(w http.ResponseWriter, _ *http.Request) {

	allDeviceInfo := m.ListAllNetDev()
	vfPod := m.fetchVfPodInfo()

	openMetContent := convertToOpnMetFormat(allDeviceInfo, vfPod)

	retResponseInOpnMetFormat(w, openMetContent)

}

// endpoint to fetch metrics related to given network
// device by name
// http://<hostname>:<port>/device/<DeviceName>
func (m *metricsHandler) deviceGet(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	DeviceName := params["DeviceName"]
	res := m.deviceByProp("Name", DeviceName)

	if res.Devices == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", DeviceName)
		return
	}

	openMetContent := convertToOpnMetFormat(res.Devices, res.VfPod)
	retResponseInOpnMetFormat(w, openMetContent)

}

// endpoint to fetch metrics related to given network
// device by pci addr
// http://<hostname>:<port>/metrics/pci-addr/<PciAddr>
func (m *metricsHandler) pciAddrGet(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	PciAddr := params["PciAddr"]
	res := m.deviceByProp("Pciaddr", PciAddr)
	var openMetContent string
	if res.Devices == nil && res.VfDev == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, " %s Device Not Found ", PciAddr)
		return
	}

	if res.VfDev != nil {
		openMetContent += vfDeviceOpenMetfmt(res.VfDev[0])
	}

	openMetContent += convertToOpnMetFormat(res.Devices, res.VfPod)
	retResponseInOpnMetFormat(w, openMetContent)

}

// function to return response from string to byte
func retResponseInOpnMetFormat(w http.ResponseWriter, cntnt string) {
	w.Header().Set("Content-Type", OpenMetContentType)
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(cntnt))

	if err != nil {
		log.Error(err)
	}
}

// function to cast PfDevices data from struct to openmetrics format
func convertToOpnMetFormat(pfd PfDevices, p VfpodInfo) string {
	// store stats data in openmetrics
	var openMetContent string

	// convert Device Information into Open metrics format
	for _, device := range pfd {
		openMetContent += devStatOpenMet(device)

		// convert virtual device info into Open metrics format
		// if present
		if len(device.VfsDetails) != 0 {
			for _, vf := range device.VfsDetails {
				openMetContent += vfDeviceOpenMetfmt(vf)
			}
		}
	}

	// convert Pods Information into Open metrics format
	for _, pod := range p {
		openMetContent += podStatsOpenMetfmt(pod)
	}

	return openMetContent
}

// Function to cast PfDevice data from struct to openmetrics format
func devStatOpenMet(devInfo PfDevice) string {
	var reg = openmetrics.NewRegistry()
	regName := "network_interface_device"

	// OpenMetrics Not showing Fields if it is empty
	// so adding a Blank string to show alias
	alias := " "
	if devInfo.Alias != "" {
		alias = devInfo.Alias
	}

	labels := []string{
		"device", "address", "broadcast", "duplex", "ifalias", "operstate", "pciaddr",
	}
	var openMetVar = regDevInfo(reg, regName, regName, labels)

	openMetVar.With(
		devInfo.Name,
		devInfo.HardwareAddr,
		devInfo.Broadcast,
		devInfo.Duplex,
		alias,
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
func regDevInfo(
	reg *openmetrics.Registry, name string, help string, labels []string,
) openmetrics.InfoFamily {

	var deviceInfo = reg.Info(openmetrics.Desc{
		Name:   name,
		Help:   help,
		Labels: labels,
	})

	return deviceInfo
}

// Function to cast Pod Details data from struct to openmetrics format
func podStatsOpenMetfmt(podDet map[string]string) string {
	var reg = openmetrics.NewRegistry()
	regName := "network_interface_vf_kube_pod"
	var openMetVar = regDevInfo(
		reg, regName, regName, []string{
			"namespace", "pod", "container", "resource", "device", "vf",
			"pciaddr", "address",
		},
	)

	openMetVar.With(
		podDet["Namespace"],
		podDet["Pod"],
		podDet["Container"],
		podDet["Resource"],
		podDet["Device"],
		podDet["Vf"],
		podDet["Pciaddr"],
		podDet["HardwareAddr"],
	)

	// create buffer to return
	var buf bytes.Buffer
	if _, err := reg.WriteTo(&buf); err != nil {
		panic(err)
	}
	// buffer.String()
	return buf.String()
}

// Function to cast VfDevice data from struct to openmetrics format
func vfDeviceOpenMetfmt(devInfo VfDevice) string {
	var reg = openmetrics.NewRegistry()
	regName := "network_interface_vf_device"
	var openMetVar = regDevInfo(
		reg,
		regName,
		regName,
		[]string{"device", "vf", "pciaddr", "address", "vlan", "spoofcheck", "trust"},
	)

	openMetVar.With(
		devInfo.PfName,
		strconv.Itoa(devInfo.VfInfo.ID),
		devInfo.Pciaddr,
		devInfo.VfInfo.Mac.String(),
		strconv.Itoa(devInfo.VfInfo.Vlan),
		BoolToOnOff(devInfo.VfInfo.Spoofchk),
		intToOnOff(devInfo.VfInfo.Trust),
	)
	statsString := vfStatsOpenMetfmt(devInfo)
	// create buffer to return
	var buf bytes.Buffer
	if _, err := reg.WriteTo(&buf); err != nil {
		panic(err)
	}
	// buffer.String()
	return buf.String() + statsString
}

// Function to cast VfDevice statistics data from struct to openmetrics format
func vfStatsOpenMetfmt(devInfo VfDevice) string {

	var reg = openmetrics.NewRegistry()

	// convert all netlink.LinkStatistics Struct to map
	fields := reflect.TypeOf(*&devInfo.VfInfo)
	// valPtr := reflect.ValueOf(devInfo.Statistics)
	values := reflect.Indirect(reflect.ValueOf(devInfo.VfInfo))
	// get the number of field for looping n times
	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)
		if strings.Contains(field.Name, "Tx") || strings.Contains(field.Name, "Rx") {
			// log.Info("Type:", field.Type, ",", field.Name, "=", value, "T==>", value.Type(), "\n")
			// register counter
			name := "network_interface_vf_" + strcase.ToSnake(field.Name)
			floatVal, _ := strconv.ParseFloat(fmt.Sprint(value), 32)
			var newInfo = reg.Counter(openmetrics.Desc{
				Name:   name,
				Help:   name,
				Labels: []string{"device", "vf", "pciaddr"},
			})
			newInfo.With(
				devInfo.PfName,
				fmt.Sprintf("%d", devInfo.VfInfo.ID),
				devInfo.Pciaddr,
			).Add(float64(floatVal))
			// log.Infof("name= %s, type=%s value=%d", name, field.Type, value.Interface().(uint32))
		}

	}
	// create buffer to return
	var buf bytes.Buffer
	if _, err := reg.WriteTo(&buf); err != nil {
		panic(err)
	}
	// buffer.String()
	return buf.String()
}
