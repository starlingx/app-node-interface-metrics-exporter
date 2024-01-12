/*
 Copyright (c) 2023-2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

// Func to set up Test environment for methods in openmet_handler
func TestInitiateOpenMetSetupSuite(t *testing.T) {
	suite.Run(t, new(OpenMetSuiteHandler))
	suite.Run(t, new(SuiteHandler))
}

// Test to verify simple metricsGet API
func (suite *OpenMetSuiteHandler) TestMetricsGet() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/metrics/", nil)
	w := httptest.NewRecorder()

	// Mocked Output
	vfPods := sampleVfPodInfo
	pfDevices := PfDevices{
		"eth0": samplePfDevWithNoStats,
	}

	// Mock Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevices)

	suite.metricsHandler.metricsGet(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

// This Test verifies the flow of fetching *eth0*
// Where eth0 doesn't have statistics
func (suite *OpenMetSuiteHandler) TestDeviceGet() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "eth0",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	pfDevice := PfDevices{
		"eth0": samplePfDevWithNoStats,
	}

	// Mock Expectations
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	// Function Call
	suite.metricsHandler.deviceGet(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

// Test verifies the 404 [Not Found] expectation from deviceGet API
func (suite *SuiteHandler) TestDeviceGetNeg() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"DeviceName": "404",
	}
	r = mux.SetURLVars(r, vars)

	// empty interface,
	// as expected, API will return 404
	expDevRes := Response{}

	// Mock Expectations
	suite.mockUtils.On("deviceByProp").Return(expDevRes)

	suite.metricsHandler.deviceGet(w, r)

	suite.Equal(http.StatusNotFound, w.Code)
}

// Test verifies the deviceGet API which has VFs
func (suite *OpenMetSuiteHandler) TestDeviceGetWithVfs() {
	// /device/{DeviceName}
	r, _ := http.NewRequest("GET", "/json/device/", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"DeviceName": "SamplePfName",
	}

	// Mocked Output
	vfPods := VfpodInfo{}
	pfDevice := PfDevices{
		"SamplePfName": SamplePfDevice,
	}

	r = mux.SetURLVars(r, vars)

	// Mock Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.deviceGet(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

// Test verifies pciAddrGet API given PciAddr as input
func (suite *OpenMetSuiteHandler) TestPciAddrGet() {
	r, _ := http.NewRequest("GET", "/pci-addr/", nil)
	w := httptest.NewRecorder()

	// Hack to try to fake gorilla/mux vars
	vars := map[string]string{
		"PciAddr": "0000:00:1f.6",
	}
	r = mux.SetURLVars(r, vars)

	// Mocked Output
	pfDevice := PfDevices{
		"eth0": samplePfDevWithNoStats,
	}

	// Mock Expectations
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)
	suite.metricsHandler.pciAddrGet(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

// Test verifies 404 [Not found] expectation from pciAddrGet API
// Test input is given as eth0, which is the device name
func (suite *OpenMetSuiteHandler) TestPciAddrGetNeg() {
	r, _ := http.NewRequest("GET", "/pci-addr/", nil)
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

	// Mock Expectations
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	suite.metricsHandler.pciAddrGet(w, r)

	suite.Equal(http.StatusNotFound, w.Code)
}

// Function Unsets the expectations set to the methods after every test
func (suite *OpenMetSuiteHandler) TearDownTest() {
	suite.mockDevHandler.On("ListAllNetDev").Unset()
	suite.mockVfPodClass.On("fetchVfPodInfo").Unset()
}
