/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMetricsGetJson(t *testing.T) {
	t.Parallel()
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/metrics/", nil)
	w := httptest.NewRecorder()

	metricsGetJSON(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
}

func TestDeviceGetJson(t *testing.T) {
	t.Parallel()
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "eth0",
	}

	r = mux.SetURLVars(r, vars)

	deviceGetJSON(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
}

func TestDeviceGetJsonNeg(t *testing.T) {
	t.Parallel()
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "404",
	}

	r = mux.SetURLVars(r, vars)

	deviceGetJSON(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPciAddrGetJsonNeg(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/json/pci-addr/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"PciAddr": "11:00:00",
	}

	r = mux.SetURLVars(r, vars)

	pciAddrGetJSON(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
