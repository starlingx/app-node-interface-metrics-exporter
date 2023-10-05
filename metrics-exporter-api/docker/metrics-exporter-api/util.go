/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"os"

	"github.com/safchain/ethtool"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// DevInfo Data structures to store Device info.
type DevInfo struct {
	Name         string
	Type         string
	HardwareAddr string
	OperState    string
	EncapType    string
	Alias        string
	Pciaddr      string
	Broadcast    string
	Statistics   *netlink.LinkStatistics
}

// OpenFile function
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

// ListAllNetDev func to findout all the Pcie Devices and its statistics
// which have network or ethernet in Classname
// func ListAllPciNetDev(sysPath string) map[string]map[string]string {
func ListAllNetDev() map[string]DevInfo {

	// make the list of dict containing DevInfo
	allDeviceInfo := make(map[string]DevInfo)

	allNetDevices, _ := netlink.LinkList()
	for _, dev := range allNetDevices {
		deviceName := dev.Attrs().Name

		pciaddr, err := ethtool.BusInfo(dev.Attrs().Name)
		if err != nil {
			log.Infof(
				"Unable to fetch Info from ethtool for %s, err: %s",
				deviceName, err,
			)
		}

		var broadcast string
		// try to fetch Broadcast Address
		addr, err := netlink.AddrList(dev, netlink.NewRule().Family)

		// we know that all device will not have Broadcast addr
		// so just logged the name of devices which are not
		if err != nil {
			log.Info("Unable to fetch AddrList ", err)
		}
		if len(addr) > 0 {
			broadcast = addr[0].Broadcast.String()
		}

		allDeviceInfo[deviceName] = DevInfo{
			Name:         deviceName,
			Type:         dev.Type(),
			HardwareAddr: dev.Attrs().HardwareAddr.String(),
			OperState:    dev.Attrs().OperState.String(),
			EncapType:    dev.Attrs().EncapType,
			Alias:        dev.Attrs().Alias,
			Pciaddr:      pciaddr,
			Broadcast:    broadcast,
			Statistics:   dev.Attrs().Statistics,
		}

	}
	return allDeviceInfo
}
