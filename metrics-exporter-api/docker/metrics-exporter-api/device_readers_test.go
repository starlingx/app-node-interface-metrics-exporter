/*
 Copyright (c) 2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vishvananda/netlink"
)

// Func to set up Test environment for methods in device_readers
func TestSetupSuiteDevReaders(t *testing.T) {
	suite.Run(t, new(DevReaderSuiteHandler))
}

// Test verifies the ListAllNetDev method flow of calls
// Asserted output by mocking Netlink Client
func (suite *DevReaderSuiteHandler) TestListAllNetDev() {
	// Test Input
	sampleNetlinkDevs := giveSampleNetlinkDev()

	// Mocked Output
	Addr := netlink.Addr{}
	AddrList := append([]netlink.Addr{}, Addr)

	// Mocked Expectations
	suite.mockNetlinkClient.On("getNetlinkDevList").Return(sampleNetlinkDevs)
	suite.mockNetlinkClient.On("getNetlinkAddrList").Return(AddrList)

	res := suite.devReceiver.ListAllNetDev()

	suite.Equal(sampleNetlinkDevs[0].Type(), res["eth0"].Type)
	suite.Equal(sampleNetlinkDevs[0].Attrs().OperState.String(), res["eth0"].OperState)
}

// Function Unsets the expectations set to the methods after every test
func (suite *DevReaderSuiteHandler) TearDownTest() {
	suite.mockNetlinkClient.On("getNetlinkDevList").Unset()
	suite.mockNetlinkClient.On("getNetlinkAddrList").Unset()
}
