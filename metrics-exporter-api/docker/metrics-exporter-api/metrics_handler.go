/*
 Copyright (c) 2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

// We add dependencies of the class MetricsHandler here
type metricsHandler struct {
	// VfPodClass is a class containing methods dealing with Vf kube Pod info
	VfPodClass

	// Class which handles connection with netlink devices
	DeviceHandler

	// Util class containing all helper methods
	Utils
}

func newMetricsHandler() *metricsHandler {
	vfPodReceiver := &VfPodReceiver{}
	devReceiver := &DevReceiver{
		NetlinkClient: &netlinkReceiver{},
	}
	utilsReceiver := &UtilReceiver{
		VfPodClass:    vfPodReceiver,
		DeviceHandler: devReceiver,
	}
	metricsHandler := &metricsHandler{
		VfPodClass:    vfPodReceiver,
		DeviceHandler: devReceiver,
		Utils:         utilsReceiver,
	}
	return metricsHandler
}
