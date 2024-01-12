/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

// Func to set up Test environment for methods in json_handler
func TestInitiateJsonSuite(t *testing.T) {
	suite.Run(t, new(SuiteHandler))
	suite.Run(t, new(JsonMetSuiteHandler))
}

// Test verifies simple metricsGetJson, empty expected output
func (suite *SuiteHandler) TestMetricsGetJson() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/metrics/", nil)
	w := httptest.NewRecorder()

	// Mocked Output - Right now we are initializing empty VfpodInfo
	vfPods := VfpodInfo{}
	pfDevices := PfDevices{}

	// Mock Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevices)

	suite.metricsHandler.metricsGetJSON(w, r)
	suite.Equal(http.StatusOK, w.Code)
	suite.Equal(w.Header().Get("Content-Type"), "application/json")
}

// Test verifies the deviceGetJson API
// Asserts on the output given when SamplePfName is the input dev name
func (suite *JsonMetSuiteHandler) TestDeviceGetJson() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "SamplePfName",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	vfPods := VfpodInfo{}
	pfDevice := PfDevices{
		"SamplePfName": SamplePfDevice,
	}

	// Mock Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.deviceGetJSON(w, r)

	var jsonResponse Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))

	suite.Equal(jsonResponse.Devices["SamplePfName"], SamplePfDevice)
	suite.Equal(err, nil)
}

// Test verifies 404 [Not found] expectation from the deviceGetJSON API
func (suite *JsonMetSuiteHandler) TestDeviceGetJsonNeg() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "404",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	vfPods := VfpodInfo{}
	pfDevice := PfDevices{
		"SamplePfName": SamplePfDevice,
	}

	// Mock Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.deviceGetJSON(w, r)
	suite.Equal(http.StatusNotFound, w.Code)
}

// Test verifies retreival of a VF dev info of pciAddrGetJSON API
// VF output is asserted with the SampleVfDevice
func (suite *JsonMetSuiteHandler) TestPciAddrGetJsonVf() {
	// t.Parallel()
	r, _ := http.NewRequest("GET", "/json/pci-addr/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"PciAddr": "SampleVfPci",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	vfPods := VfpodInfo{}
	pfDevice := PfDevices{
		"SamplePfName": SamplePfDevice,
	}

	// Mocked Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.pciAddrGetJSON(w, r)

	var jsonResponse Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &jsonResponse)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("application/json", w.Header().Get("Content-Type"))
	suite.Equal(jsonResponse.VfDev[0], SampleVfDevice)
	suite.Equal(err, nil)
}

// Test verifies 404 [Not found] expectation from pciAddrGetJSON API
// which is to give empty response even if the device is present
// in the Net Dev List but dev name is given as input
func (suite *JsonMetSuiteHandler) TestPciAddrGetJsonNeg() {
	r, _ := http.NewRequest("GET", "/json/pci-addr/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	// Passing deviceName to PciAddrGet method which is expected to fail
	vars := map[string]string{
		"PciAddr": "eth0",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	pfDevice := PfDevices{
		"eth0": samplePfDevWithNoStats,
	}

	// Mocked Expectation
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.pciAddrGetJSON(w, r)

	// suite.mockVfPodClass.AssertNotCalled(suite.T(), "fetchVfPodInfo")
	suite.Equal(http.StatusNotFound, w.Code)
}

// Function Unsets the expectations set to the methods after every test
func (suite *JsonMetSuiteHandler) TearDownTest() {
	suite.mockDevHandler.On("ListAllNetDev").Unset()
	suite.mockVfPodClass.On("fetchVfPodInfo").Unset()
}
