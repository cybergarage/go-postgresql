// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
go-postgresqld is an example of a compatible PostgreSQL server implementation using go-postgresql.

	NAME
	 go-postgresqld

	SYNOPSIS
	 go-postgresqld [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	clog "github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-postgresql/examples/go-postgresqld/server"
)

const (
	ProgramName = " go-postgresqld"
)

func main() {
	isDebugEnabled := flag.Bool("debug", false, "enable debugging log output")
	isTraceEnabled := flag.Bool("trace", false, "enable trace log output")
	isProfileEnabled := flag.Bool("profile", false, "enable profiling server")
	flag.Parse()

	logLevel := clog.LevelInfo
	if *isTraceEnabled {
		logLevel = clog.LevelTrace
	}
	if *isDebugEnabled {
		logLevel = clog.LevelDebug
	}
	clog.SetSharedLogger(clog.NewStdoutLogger(logLevel))

	if *isProfileEnabled {
		go func() {
			// nolint: gosec
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Start server

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		log.Printf("%s couldn't be started (%s)", ProgramName, err.Error())
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				log.Printf("Caught SIGHUP, restarting...")
				err = server.Restart()
				if err != nil {
					log.Printf("%s couldn't be restarted (%s)", ProgramName, err.Error())
					os.Exit(1)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				log.Printf("Caught %s, stopping...", s.String())
				err = server.Stop()
				if err != nil {
					log.Printf("%s couldn't be stopped (%s)", ProgramName, err.Error())
					os.Exit(1)
				}
				exitCh <- 0
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)
}
