/*
Copyright (c) 2023-2024 Wind River Systems, Inc.

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

// Class for all Device Readers Method
type DeviceHandler interface {
	ListAllNetDev() PfDevices
}

// This struct encapsulates the dependencies of device_readers methods
// aids in creating mocks for these dependencies and test functions
type DevReceiver struct {
	NetlinkClient
}

// Mapper Function
type Mapper func(pfd *PfDevice) string

// Mapper to retrieve the given property as a string
// for PfDevice.
var PropMapper = map[string]Mapper{
	"Name":    func(pfd *PfDevice) string { return pfd.Name },
	"Pciaddr": func(pfd *PfDevice) string { return pfd.Pciaddr },
}

// PfDevices structure
type PfDevices map[string]PfDevice
type vfDevices []VfDevice

// Function to list all the network devices that are not of type
// ‘veth’ and are in the UP state.
func (dr *DevReceiver) ListAllNetDev() PfDevices {

	allDeviceInfo := PfDevices{}
	// fetch all devices from Netlink Library
	netLinkDevices, _ := dr.getNetlinkDevList()

	for _, device := range netLinkDevices {
		// ignore veth devices
		if device.Type() != "veth" {

			deviceName := device.Attrs().Name

			// Because the majority of the data for downed devices is
			// missing, only retrieve up-and-running devices.
			if device.Attrs().OperState.String() == "up" {
				// try to fetch Broadcast Address
				addr, _ := dr.getNetlinkAddrList(device, netlink.NewRule().Family)

				// Check if devices has vfs then add
				vfs := vfDevices{}
				if len(device.Attrs().Vfs) != 0 {
					vfs, _ = dr.updateVirtualDev(device)
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
					Broadcast:    fmt.Sprint(addr),
					Duplex:       intToDuplex(fetchEthToolData(deviceName, "Duplex")),
					VfsDetails:   vfs,
					Statistics:   device.Attrs().Statistics,
				}
			}

		}
	}
	return allDeviceInfo
}

func (dr *DevReceiver) updateVirtualDev(d netlink.Link) (vfDevices, error) {
	vfs := vfDevices{}
	// Read all VF directories and add the VF PCI address to the vfList.
	// For example, the directory path for the VF can be
	// /sys/class/net/enp129s0f1/device/virtfn*<num>.
	for _, vf := range dr.getVfsDirNames(d.Attrs().Name) {

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
func (dr *DevReceiver) getVfsDirNames(deviceName string) []string {
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

// Interface to hold netlink client calls
type NetlinkClient interface {
	getNetlinkDevList() ([]netlink.Link, error)
	getNetlinkAddrList(device netlink.Link, family int) ([]netlink.Addr, error)
}

// This struct encapsulates the dependencies of netlink methods
// aids in creating mocks for these dependencies and test functions
type netlinkReceiver struct {
	// empty
}

// func to fetch All Netlink Devices
func (n *netlinkReceiver) getNetlinkDevList() ([]netlink.Link, error) {
	return netlink.LinkList()
}

// func to fetch Broadcast addr for a dev
func (n *netlinkReceiver) getNetlinkAddrList(device netlink.Link, family int) ([]netlink.Addr, error) {
	return netlink.AddrList(device, family)
}
