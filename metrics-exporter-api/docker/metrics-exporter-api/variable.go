/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"flag"
	"time"
)

var (
	addr = flag.String(
		"web.listen-address", ":9110", "Port to listen on for web interface.",
	)
	logLevel = flag.String(
		"log.level", "info", "log level. Valid options trace,"+
			" debug, info, warning, error, fatal and panic",
	)
	logFileName = flag.String(
		"log.file", "node_metrics_api.log", "Log file name",
	)

	// StartupTime string
	StartupTime int64

	// OpenMetContentType content Type
	// OpenMetContentType = "application/openmetrics-text; version=1.0.0; charset=utf-8"
	OpenMetContentType = "text/plain"

	// Version support
	Version = "Dev"
	//Build datetime
	Build = time.Now().String()
)
