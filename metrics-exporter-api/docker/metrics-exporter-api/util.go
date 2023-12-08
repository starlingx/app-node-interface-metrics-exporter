/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
)

// Class Def for all Utility Methods
type Utils interface {
	deviceByProp(propName string, propVal string) Response
	getMatchedPod(vfs vfDevices) VfpodInfo
	findPciAddrInVF(pciAddr string) Response
}

// This struct encapsulates the dependencies of util methods
// aids in creating mocks for these dependencies and test functions
type UtilReceiver struct {
	DeviceHandler

	VfPodClass
}

// OpenFile function to open a file with the given file path.
func OpenFile(fileName string) *os.File {
	logFile := fileName
	// Open logfile
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Info("Failed to create logfile" + logFile)
		log.Error(err)
		panic(err)
	}
	return f
}

// it is required to get the annotations from pod description
// func to cast string to map
func castStr2Map(in string) []map[string]interface{} {
	var jsonMap []map[string]interface{}
	err := json.Unmarshal([]byte(in), &jsonMap)
	if err != nil {
		log.Info("Unable to marshal ==>", in)
	}
	return jsonMap
}

// func to get value from map keys
func getValfromInterface(val interface{}, key string) interface{} {
	iter := reflect.ValueOf(val).MapRange()
	for iter.Next() {
		// key := iter.Key().Interface()
		// value := iter.Value().Interface()
		if iter.Key().Interface() == key {
			return iter.Value().Interface()
		}
	}
	return nil
}

// EvalSymlinks Required to enable testing (filepath.EvalSymlinks does not
// support the fs.FS interface that fstest implements)
var EvalSymlinks = func(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

// func to get pod details if pciaddr is matched in any virtual devices
// pod have pciaddr of only virtual devices
func (u *UtilReceiver) getMatchedPod(vfs vfDevices) VfpodInfo {
	// virtual function Pod Info
	vfPod := u.fetchVfPodInfo()
	matchedVfPod := VfpodInfo{}

	for _, vf := range vfs {
		// try to match the pci address from pf->vfs to vf assigned to pod
		for key, pod := range vfPod {
			if pod["Pciaddr"] == vf.Pciaddr {
				matchedVfPod[key] = pod
			}
		}
	}
	return matchedVfPod
}

// func to find device by its proerty name
// e.g: find name:en03, find pciAddr: 0000:81:0a.0
func (u *UtilReceiver) deviceByProp(propName string, propVal string) Response {
	// deviceRec := new(DevReceiver)
	devStats := u.findByProperty(propName, propVal)

	if len(devStats) == 0 {
		// before failing search pci addr in vfs
		// if the propName is Pciaddr
		// the look into virtual devices as well
		if propName == "Pciaddr" {
			log.Info("Searching in VFS")
			return u.findPciAddrInVF(propVal)
		}

		return Response{}
	}

	var matchedPod = VfpodInfo{}
	// if found
	if len(devStats) != 0 {
		// check if vfs exists
		if len(devStats[propVal].VfsDetails) != 0 {
			// check any of the vfs assigned to pod
			matchedPod = u.getMatchedPod(devStats[propVal].VfsDetails)
		}
	}

	return Response{Devices: devStats, VfPod: matchedPod}
}

// func to find PciAddr in virtual function
func (u *UtilReceiver) findPciAddrInVF(pciAddr string) Response {
	allPf := u.ListAllNetDev()
	var vfDevices []VfDevice

	for _, devStats := range allPf {
		if len(devStats.VfsDetails) != 0 {
			log.Info("Search here in VFS")
			for _, vf := range devStats.VfsDetails {
				if vf.Pciaddr == pciAddr {
					vfDevices = append(vfDevices, vf)
					return Response{VfDev: vfDevices, VfPod: u.getMatchedPod(vfDevices)}
				}
			}
		}
	}
	return Response{}

}

// A getter function to retrieve the given property as a string
// for PfDevice.
func (pfd *PfDevice) getProperty(property string) string {
	return PropMapper[property](pfd)
}

// Function to find a device by a given property name and its value
// for PfDevice.
func (u *UtilReceiver) findByProperty(propName string, propVal string) PfDevices {
	allDeviceStat := u.ListAllNetDev()
	var PfDevices = make(map[string]PfDevice)

	for _, devStats := range allDeviceStat {
		if devStats.getProperty(propName) == propVal {
			PfDevices[propVal] = devStats
		}
	}
	return PfDevices
}

// intToString function to cast int to string
// 0 is off , 1 is on else blank
func intToOnOff(t uint32) string {
	var retVal string
	switch t {
	case 0:
		retVal = "off"
	case 1:
		retVal = "on"
	default:
		retVal = ""
	}
	return retVal
}

// BoolToOnOff function to cast Bool to string
// if t is 0 then off else always on
func BoolToOnOff(t bool) string {
	retVal := "off"
	if t {
		retVal = "on"
	}
	return retVal
}

// intToDuplex function to cast uint to string
//
//	0 is half, 1 is full , others is blank
func intToDuplex(d int) string {
	var retVal string
	switch d {
	case 0:
		retVal = "half"
	case 1:
		retVal = "full"
	default:
		retVal = ""
	}
	return retVal
}

// fetchPciAddr function to retrieve the pciaddr from the ethtool
// library for a specified device
func fetchPciAddr(devName string) string {
	pciaddr, err := ethtool.BusInfo(devName)

	if err != nil {
		log.Warnf(
			"Unable to fetch pciaddr from ethtool for %s, err: %s",
			devName, err,
		)
	}
	return pciaddr
}

// fetchEthToolData function to retrieve the information from the ethtool
// library for a specified device
func fetchEthToolData(devName string, key string) int {
	// get Duplex information from ethtool
	ethinfo, err := ethtool.CmdGetMapped(devName)
	if err != nil {
		log.Warnf(
			"Unable to fetch ethinfo from ethtool for %s, err: %s",
			devName, err,
		)
	}
	return int(ethinfo[key])
}
