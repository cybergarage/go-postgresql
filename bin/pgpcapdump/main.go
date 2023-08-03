// Copyright (C) 2019 Satoshi Konno. All rights reserved.
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
pgpcapdump is a dump utility for PostgreSQL packet capture file.

	NAME
	 pgpcapdump

	SYNOPSIS
	 pgpcapdump [OPTIONS] FILE

	OPTIONS
	-q      : Extract queries.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
//nolint:forbidigo
package main

import (
	"bufio"
	"bytes"
	"flag"
	"os"

	"github.com/cybergarage/go-postgresql/postgresql/protocol/message"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	ProgramName = "pgpcapdump"
)

func usages() {
	println("Usage:")
	println("  " + ProgramName + " [OPTIONS] FILE")
	println("")
	println("Options:")
	println("  -q      : Extract queries.")
	println("")
	println("Return Value:")
	println("  Return EXIT_SUCCESS or EXIT_FAILURE")
	os.Exit(1)
}

func exit(err error) {
	println("Error: " + err.Error())
	os.Exit(1)
}

func main() {
	isQueryEnabled := flag.Bool("q", true, "extract queries")
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		usages()
	}

	pcapFilename := args[0]
	handle, err := pcap.OpenOffline(pcapFilename)
	if err != nil {
		exit(err)
	}

	pktSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for pkt := range pktSource.Packets() {
		reader := message.NewMessageReaderWith(bufio.NewReader(bytes.NewReader(pkt.Data())))
		msg, err := message.NewMessageWithReader(reader)
		if err != nil {
			exit(err)
		}
		reader = message.NewMessageReaderWith(bufio.NewReader(bytes.NewReader(pkt.Data())))
		switch msg.Type { // nolint:exhaustive
		case message.QueryMessage:
			if !*isQueryEnabled {
				query, err := message.NewQueryWithReader(reader)
				if err != nil {
					exit(err)
				}
				println(query.Query)
			}
		}
	}

	os.Exit(0)
}
