/*
 Copyright (c) 2023 Wind River Systems, Inc.

 SPDX-License-Identifier: Apache-2.0

 All Rights Reserved.
*/

package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	handler := allHandlers()

	// StartUp Time used to calculate Uptime for the App
	atomic.StoreInt64(&StartupTime, time.Now().UnixNano())

	// Parse the flags or commandline arguments
	flag.Parse()

	// log setup
	logFile := OpenFile(*logFileName)
	// close file on exit
	defer logFile.Close()
	logOutput := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logOutput)
	logLevel, _ := log.ParseLevel(*logLevel)
	// Only log the debug severity or above
	log.SetLevel(logLevel)

	// Print all flags what valuse is used
	flag.VisitAll(func(f *flag.Flag) {
		log.Infof("%s: %s", f.Name, f.Value)
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, handler)
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           loggedRouter,
	}
	// Starting Http server
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	// on exit log
	log.Info("Server stopped\n")
}
