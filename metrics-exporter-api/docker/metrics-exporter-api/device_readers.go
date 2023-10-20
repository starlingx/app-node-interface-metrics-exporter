/*
Copyright (c) 2023 Wind River Systems, Inc.

SPDX-License-Identifier: Apache-2.0

All Rights Reserved.
*/

package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// PfDevice Data structures to store Device info.
type PfDevice struct {
	Name         string
	Type         string
	DevType      string
	HardwareAddr string
	OperState    string
	EncapType    string
	Alias        string
	Pciaddr      string
	Broadcast    string
	Duplex       string
	VfsDetails   vfDevices
	Statistics   *netlink.LinkStatistics
}

// VfDevice Data structure to store VF Info
type VfDevice struct {
	PfName  string
	Pciaddr string
	Address string
	VfInfo  netlink.VfInfo
}

// Mapper Function
type Mapper func(pfd *PfDevice) string

// PropMapper Function
var PropMapper = map[string]Mapper{
	"Name":    func(pfd *PfDevice) string { return pfd.Name },
	"Pciaddr": func(pfd *PfDevice) string { return pfd.Pciaddr },
}

// PfDevices structure
type PfDevices map[string]PfDevice
type vfDevices []VfDevice

// A getter function to retrieve the given property as a string
// for PfDevice.
func (pfd *PfDevice) getProperty(property string) string {
	return PropMapper[property](pfd)
}

// Function to find a device by a given property name and its value
// for PfDevice.
func (pfd PfDevices) findByProperty(propName string, propVal string) PfDevices {
	allDeviceStat := ListAllNetDev()
	var PfDevices = make(map[string]PfDevice)

	for _, devStats := range allDeviceStat {
		if devStats.getProperty(propName) == propVal {
			PfDevices[propVal] = devStats
		}
	}
	return PfDevices
}

// Function to list all the network devices that are not of type
// ‘veth’ and are in the UP state.
func ListAllNetDev() PfDevices {

	allDeviceInfo := PfDevices{}
	// fetch all devices from Netlink Library
	netLinkDevices, _ := netlink.LinkList()

	for _, device := range netLinkDevices {
		// ignore veth devices
		if device.Type() != "veth" {

			deviceName := device.Attrs().Name

			// Because the majority of the data for downed devices is
			// missing, only retrieve up-and-running devices.
			if device.Attrs().OperState.String() == "up" {
				// try to fetch Broadcast Address
				addr, _ := netlink.AddrList(device, netlink.NewRule().Family)

				// Check if devices has vfs then add
				vfs := vfDevices{}
				if len(device.Attrs().Vfs) != 0 {
					vfs, _ = updateVirtualDev(device)
				}

				allDeviceInfo[deviceName] = PfDevice{
					Name:         deviceName,
					Type:         device.Type(),
					DevType:      "PHYSICAL",
					HardwareAddr: device.Attrs().HardwareAddr.String(),
					OperState:    device.Attrs().OperState.String(),
					EncapType:    device.Attrs().EncapType,
					Alias:        device.Attrs().Alias,
					Pciaddr:      fetchPciAddr(deviceName),
					Broadcast:    fmt.Sprint(addr[0]),
					Duplex:       intToDuplex(fetchEthToolData(deviceName, "Duplex")),
					VfsDetails:   vfs,
					Statistics:   device.Attrs().Statistics,
				}
			}

		}
	}
	return allDeviceInfo
}

func updateVirtualDev(d netlink.Link) (vfDevices, error) {
	vfs := vfDevices{}
	// Read all VF directories and add the VF PCI address to the vfList.
	// For example, the directory path for the VF can be
	// /sys/class/net/enp129s0f1/device/virtfn*<num>.
	for _, vf := range getVfsDirNames(d.Attrs().Name) {

		if link, err := EvalSymlinks(vf); err == nil {
			vfID, _ := strconv.Atoi(filepath.Base(vf)[6:])

			// we can have mac to long notation by printing
			// log.Infof(" %+s", vf_info.Mac.String())
			vfs = append(vfs, VfDevice{
				VfInfo:  d.Attrs().Vfs[vfID],
				Pciaddr: filepath.Base(link),
				Address: d.Attrs().Vfs[vfID].Mac.String(),
				PfName:  d.Attrs().Name,
			})
			// we can have mac to long notation by printing
			// log.Infof(" %+s", vfs[vfID].VfInfo.Mac.String())

		} else {
			log.Printf("error evaluating symlink '%s'\n%v", vf, err)
			return vfs, err
		}

	}
	return vfs, nil
}

// getVfsDirNames function retrieves all the VF directories from the
// file system. It returns the folders inside the PF that have the
// pattern ‘virtfn*<num>’. For example, the directory path for the VF can be
// /sys/class/net/enp129s0f1/device/virtfn*<num>.
func getVfsDirNames(deviceName string) []string {
	// create a path to get all virtual dev
	// e.g, sys/class/net/enp129s0f1/device
	devicePath := filepath.Join(
		*sysPath, "/class/net/", deviceName, "/device/virtfn*",
	)
	// find all the vfs
	vfsFile, err := filepath.Glob(devicePath)

	if err != nil {
		log.Warnf("Invalid pattern\n%v", err) // unreachable code
	}
	return vfsFile
}
