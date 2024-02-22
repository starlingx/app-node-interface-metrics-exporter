/*
 Copyright (c) 2024 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/vishvananda/netlink"
)

type SuiteHandler struct {
	suite.Suite
	mockVfPodClass *MockVfPodClass
	mockDevHandler *MockDeviceHandler
	mockUtils      *MockUtils
	metricsHandler *metricsHandler
}

func (suite *SuiteHandler) SetupSuite() {
	// Common setup logic
	suite.mockVfPodClass = new(MockVfPodClass)
	suite.mockDevHandler = new(MockDeviceHandler)
	suite.mockUtils = new(MockUtils)
	suite.metricsHandler = &metricsHandler{
		VfPodClass:    suite.mockVfPodClass,
		DeviceHandler: suite.mockDevHandler,
		Utils:         suite.mockUtils,
	}
}

type OpenMetSuiteHandler struct {
	suite.Suite
	mockVfPodClass *MockVfPodClass
	mockDevHandler *MockDeviceHandler
	mockUtils      Utils
	metricsHandler *metricsHandler
}

func (suite *OpenMetSuiteHandler) SetupSuite() {
	suite.mockVfPodClass = new(MockVfPodClass)
	suite.mockDevHandler = new(MockDeviceHandler)
	suite.mockUtils = &UtilReceiver{
		DeviceHandler: suite.mockDevHandler,
		VfPodClass:    suite.mockVfPodClass,
	}
	suite.metricsHandler = &metricsHandler{
		VfPodClass:    suite.mockVfPodClass,
		DeviceHandler: suite.mockDevHandler,
		Utils:         suite.mockUtils,
	}
}

type JsonMetSuiteHandler struct {
	suite.Suite
	mockVfPodClass *MockVfPodClass
	mockDevHandler *MockDeviceHandler
	mockUtils      Utils
	metricsHandler *metricsHandler
}

func (suite *JsonMetSuiteHandler) SetupSuite() {
	suite.mockVfPodClass = new(MockVfPodClass)
	suite.mockDevHandler = new(MockDeviceHandler)
	suite.mockUtils = &UtilReceiver{
		DeviceHandler: suite.mockDevHandler,
		VfPodClass:    suite.mockVfPodClass,
	}
	suite.metricsHandler = &metricsHandler{
		VfPodClass:    suite.mockVfPodClass,
		DeviceHandler: suite.mockDevHandler,
		Utils:         suite.mockUtils,
	}
}

type UtilSuiteHandler struct {
	suite.Suite
	mockVfPodClass *MockVfPodClass
	mockDevHandler *MockDeviceHandler
	UtilReceiver   *UtilReceiver
}

func (suite *UtilSuiteHandler) SetupSuite() {
	suite.mockVfPodClass = new(MockVfPodClass)
	suite.mockDevHandler = new(MockDeviceHandler)
	suite.UtilReceiver = &UtilReceiver{
		DeviceHandler: suite.mockDevHandler,
		VfPodClass:    suite.mockVfPodClass,
	}
}

type DevReaderSuiteHandler struct {
	suite.Suite
	mockNetlinkClient *MockNetlinkClient
	devReceiver       *DevReceiver
}

func (suite *DevReaderSuiteHandler) SetupSuite() {
	suite.mockNetlinkClient = new(MockNetlinkClient)
	suite.devReceiver = &DevReceiver{
		NetlinkClient: suite.mockNetlinkClient,
	}
}

// Code partition for mocking types
type MockVfPodClass struct {
	mock.Mock
}

type MockDeviceHandler struct {
	mock.Mock
}

type MockUtils struct {
	mock.Mock
}

type MockNetlinkClient struct {
	mock.Mock
}

// Code partition for mocking implementations
func (m *MockVfPodClass) fetchVfPodInfo() VfpodInfo {
	args := m.Called()
	return args.Get(0).(VfpodInfo)
}

func (m *MockDeviceHandler) ListAllNetDev() PfDevices {
	args := m.Called()
	return args.Get(0).(PfDevices)
}

func (m *MockUtils) deviceByProp(propName string, propVal string) Response {
	args := m.Called()
	return args.Get(0).(Response)
}

func (m *MockUtils) getMatchedPod(vfs vfDevices) VfpodInfo {
	args := m.Called()
	return args.Get(0).(VfpodInfo)
}

func (m *MockUtils) findPciAddrInVF(pciAddr string) Response {
	args := m.Called()
	return args.Get(0).(Response)
}

func (m *MockNetlinkClient) getNetlinkDevList() ([]netlink.Link, error) {
	args := m.Called()
	return args.Get(0).([]netlink.Link), nil
}

func (m *MockNetlinkClient) getNetlinkAddrList(
	device netlink.Link, family int,
) []netlink.Addr {
	args := m.Called()
	return args.Get(0).([]netlink.Addr)
}

var SamplePfDevice = PfDevice{
	Name:         "SamplePfName",
	Type:         "device",
	DevType:      "PHYSICAL",
	HardwareAddr: "SampleHardwareAddr",
	OperState:    "up",
	EncapType:    "ether",
	Alias:        "",
	Pciaddr:      "SamplePciAddr",
	Broadcast:    "",
	Duplex:       "full",
	Statistics:   &netlink.LinkStatistics{},
	VfsDetails:   append(vfDevices{}, SampleVfDevice),
}

var SampleVfDevice = VfDevice{
	PfName:  "SamplePfName",
	Pciaddr: "SampleVfPci",
	Address: "SampleVfAddress",
	VfInfo:  netlink.VfInfo{},
}

var samplePfDevWithNoStats = PfDevice{
	Name:         "eth0",
	Type:         "device",
	DevType:      "PHYSICAL",
	HardwareAddr: "cc:96:e5:22:f5:a9",
	OperState:    "up",
	EncapType:    "ether",
	Alias:        "",
	Pciaddr:      "0000:00:1f.6",
	Broadcast:    "128.224.67.91/24 eth0",
	Duplex:       "full",
	Statistics:   &netlink.LinkStatistics{},
}

var sampleVfPodInfo = VfpodInfo{
	"SampleContainer-0": map[string]string{
		"Namespace":    "default",
		"Pod":          "SomePodName",
		"Container":    "SampleContainer",
		"Pciaddr":      "SampleVfPci",
		"HardwareAddr": "SampleVfAddress",
		"Resource":     "",
		"Device":       "",
		"Vf":           "1",
	},
}

func giveSampleNetlinkDev() []netlink.Link {
	netlinkDev := []netlink.Link{}
	eth0Link := &netlink.Dummy{
		LinkAttrs: netlink.LinkAttrs{
			Index:     0,
			MTU:       1500,
			Name:      "eth0",
			OperState: netlink.LinkOperState(netlink.OperUp),
		},
	}
	netlinkDev = append(netlinkDev, eth0Link)
	return netlinkDev
}
