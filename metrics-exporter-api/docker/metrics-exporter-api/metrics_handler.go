/*
 Copyright (c) 2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// We add dependencies of the class MetricsHandler here
type metricsHandler struct {
	// VfPodClass is a class containing methods dealing with Vf kube Pod info
	VfPodClass

	// Class which handles connection with netlink devices
	DeviceHandler

	// Util class containing all helper methods
	Utils

	// Class to Handle Unexpected Panics
	PanicReceiver
}

func newMetricsHandler() *metricsHandler {
	vfPodReceiver := &VfPodReceiver{}
	panicReceiver := &PanicReceiver{}
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
		PanicReceiver: *panicReceiver,
	}
	return metricsHandler
}

type PanicReceiver struct {
	// empty struct
}

func (r *PanicReceiver) recoverFromPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Error("Encountered an unexpected error: ", r)

		stack := make([]byte, 4096)
		length := runtime.Stack(stack, true)
		log.Error("Stack Trace:\n", string(stack[:length]))
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, fmt.Sprintf("%s", r), http.StatusInternalServerError)
	}
}
