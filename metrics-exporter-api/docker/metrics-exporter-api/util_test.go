/*
 Copyright (c) 2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Func to set up Test environment for methods in util.go
func TestInitiateUtilSuite(t *testing.T) {
	suite.Run(t, new(UtilSuiteHandler))
}

// Test verifies getMatchedPod method
// Asserted based on the output returned by fetchVfPodInfo method
func (suite *UtilSuiteHandler) TestGetMatchedPod() {
	// Test Input
	vfs := append(vfDevices{}, SampleVfDevice)

	// Mocked Output
	vfPods := sampleVfPodInfo

	// Mocked Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)

	res := suite.UtilReceiver.getMatchedPod(vfs)
	suite.Equal(sampleVfPodInfo, res)
}

// Test verifies flow and output of deviceByProp method
// for propName = "Name"
func (suite *UtilSuiteHandler) TestDeviceByProp() {
	// Test Input
	propName := "Name"
	propVal := "SamplePfName"

	// Mocked Output
	vfPods := VfpodInfo{}
	pfDevice := PfDevices{
		"SamplePfName": SamplePfDevice,
	}

	// Mocked Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevice)

	res := suite.UtilReceiver.deviceByProp(propName, propVal)
	suite.Equal(pfDevice, res.Devices)
}

// Test verifies if getMatchedPod returns empty when fetchVfPodInfo method
// doesn't find any VF assosciated to Kube Pods
func (suite *UtilSuiteHandler) TestGetMatchedPodNeg() {

	// Test Input
	vfs := append(vfDevices{}, SampleVfDevice)

	// Mocked Output
	vfPods := VfpodInfo{}

	// Mocked Expectations
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)

	res := suite.UtilReceiver.getMatchedPod(vfs)

	// VfPodInfo is empty, So getMatchedPod is expected to return nothing
	suite.Empty(res)
}

// Test verifies the flow and output of method findPciAddrInVF
// When a VF Pci addr is given as input, This method is supposed to search
// in all PF details to find the VF info
func (suite *UtilSuiteHandler) TestFindPciAddrInVf() {
	// Test Input
	pciAddr := "SampleVfPci"

	// Mocked Output
	pfDevices := PfDevices{
		"SamplePfName": SamplePfDevice,
	}
	vfPods := sampleVfPodInfo

	// Mocked Expectations
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevices)
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)

	res := suite.UtilReceiver.findPciAddrInVF(pciAddr)

	// suite.mockVfPodClass.AssertNumberOfCalls(suite.T(), "fetchVfPodInfo", 1)
	suite.Equal(SamplePfDevice.VfsDetails, res.VfDev)
	suite.Equal(vfPods, res.VfPod)
}

// Test verifies the flow and empty response expectation from findPciAddrInVF
// If given PCI addr is not of any VF, func is expected to return nothing
func (suite *UtilSuiteHandler) TestFindPciAddrInVfNeg() {

	// Test Input, Some random PCI addr
	pciAddr := "11:00:00"

	// Mocked Output
	pfDevices := PfDevices{
		"SamplePfName": SamplePfDevice,
	}
	vfPods := sampleVfPodInfo

	// Mocked Expectations
	suite.mockDevHandler.On("ListAllNetDev").Return(pfDevices)
	suite.mockVfPodClass.On("fetchVfPodInfo").Return(vfPods)

	res := suite.UtilReceiver.findPciAddrInVF(pciAddr)

	// suite.mockVfPodClass.AssertNumberOfCalls(suite.T(), "fetchVfPodInfo", 0)
	suite.Empty(res)
}

// Function Unsets the expectations set to the methods after every test
func (suite *UtilSuiteHandler) TearDownTest() {
	suite.mockVfPodClass.On("fetchVfPodInfo").Unset()
	suite.mockDevHandler.On("ListAllNetDev").Unset()
}
